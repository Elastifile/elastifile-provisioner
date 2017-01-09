package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"

	"github.com/elastifile/emanage-go/pkg/ejson"
	"github.com/elastifile/emanage-go/pkg/retry"

	"github.com/elastifile/errors"
)

var Log = log.New("package", "rest")

func init() {
	Log.SetHandler(log.DiscardHandler())
}

type HttpMethod string

const (
	MethodPost   HttpMethod = "POST"
	MethodPut    HttpMethod = "PUT"
	MethodGet    HttpMethod = "GET"
	MethodDelete HttpMethod = "DELETE"
)

const sessionsUri = "api/sessions"

var securityPrefix = []byte(")]}',\n")

var Timeout time.Duration = 2 * time.Hour
var DumpHTTP bool
var DumpHTTPOnError = true

var AfterShutdown func()
var BeforeStart func()
var BeforeForceReset func()

type Session struct {
	baseURL     *url.URL // Base URL of server to connect to, e.g. http://func11-cm/
	client      http.Client
	credentials credentials
	cookies     []*http.Cookie
	xsrf        string
}

func NewSession(baseURL *url.URL) *Session {
	result := &Session{
		baseURL: baseURL,
	}
	result.init()
	return result
}

func (rs *Session) init() {
	rs.client = http.Client{
		Transport: &http.Transport{DisableKeepAlives: true},
	}
}

type credentials struct {
	User     string `json:"login"`
	Password string `json:"password"`
}

func (rs *Session) Login(user string, password string) error {
	creds := credentials{user, password}
	params := struct {
		User credentials `json:"user"`
	}{creds}
	rs.credentials = creds

	jsonBody, stdErr := json.Marshal(params)
	if stdErr != nil {
		return errors.New(stdErr)
	}

	resp, _, err := rs.requestHttp(MethodPost, sessionsUri, jsonBody)
	if err != nil {
		return err
	}

	rs.cookies = resp.Cookies()
	xsrf, err := func() (string, error) {
		for _, cookie := range rs.cookies {
			if cookie.Name == "XSRF-TOKEN" {
				xsrf, e := url.QueryUnescape(cookie.Value)
				if e != nil {
					return "", e
				}

				return xsrf, nil
			}
		}
		return "", NewRestError("XSRF cookie not found", resp, []byte(""))
	}()
	if err != nil {
		return errors.New(err)
	}
	rs.xsrf = xsrf

	return nil
}

func (rs *Session) Logout() error {
	if err := rs.Request(MethodDelete, sessionsUri, nil, nil); err != nil {
		return err
	}

	rs.cookies = nil
	rs.xsrf = ""
	return nil
}

func (rs *Session) requestHttp(method HttpMethod, relURL string, body []byte) (resp *http.Response, resBody []byte, resErr error) {
	fullURL := fmt.Sprintf("%s/%s", rs.baseURL, relURL)

	req, err := http.NewRequest(string(method), fullURL, bytes.NewReader(body))
	if err != nil {
		resErr = errors.New(err)
		return
	}

	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	if rs.cookies != nil {
		for _, cookie := range rs.cookies {
			req.AddCookie(cookie)
		}
		req.Header.Add("X-XSRF-TOKEN", rs.xsrf)
	}

	if DumpHTTP {
		e := dumpRequest(req)
		if e != nil {
			resErr = errors.New(e)
			return
		}
	}

	rs.client.Timeout = Timeout
	resp, err = rs.client.Do(req)
	if err != nil {
		if DumpHTTPOnError {
			_ = dumpRequest(req)
		}
		if resp != nil && DumpHTTPOnError {
			_ = dumpResponse(resp)
		}
		resErr = err
		return
	}

	if DumpHTTP {
		e := dumpResponse(resp)
		if e != nil {
			resErr = e
			return
		}
	}

	shouldClose := true
	defer func() {
		if !shouldClose {
			return
		}
		e := resp.Body.Close()
		if e != nil && resErr == nil {
			resErr = e
		}
	}()

	resBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		if _, ok := err.(*net.OpError); ok {
			// Connection may be closed by server on logout; ignore.
			if method == MethodDelete && relURL == sessionsUri {
				shouldClose = false
				return
			}
		}
		Log.Error("Error while reading response body", "resp.Body", resp.Body)
		resErr = err
		return
	}

	if resp.StatusCode >= http.StatusBadRequest {
		resErr = NewRestError("HTTP request failed", resp, resBody)
		return
	}

	return
}

func dumpRequest(req *http.Request) error {
	data, e := httputil.DumpRequest(req, true)
	if e != nil {
		return e
	}
	fmt.Println("## HTTP DUMP REQUEST ##")
	fmt.Printf("%v\n", string(data))
	return nil
}

func dumpResponse(resp *http.Response) error {
	data, e := httputil.DumpResponse(resp, true)
	if e != nil {
		return e
	}
	fmt.Println("## HTTP DUMP RESPONSE ##")
	fmt.Printf("%v\n", string(data))
	return nil
}

//go:generate stringer -type=ControlTaskStatus

type ControlTaskStatus int

type controlTask struct {
	Status    ControlTaskStatus `json:"status"`
	LastError string            `json:"last_error"`
}

const (
	ControlTaskStatusSuccess    ControlTaskStatus = 0
	ControlTaskStatusError      ControlTaskStatus = 1
	ControlTaskStatusCanceled   ControlTaskStatus = 2
	ControlTaskStatusIncomplete ControlTaskStatus = 3
)

var taskStatuses = map[string]ControlTaskStatus{
	"success":     ControlTaskStatusSuccess,
	"error":       ControlTaskStatusError,
	"canceled":    ControlTaskStatusCanceled,
	"in_progress": ControlTaskStatusIncomplete,
}

func (cts *ControlTaskStatus) UnmarshalJSON(data []byte) (err error) {
	var (
		intVal int
		strVal string
	)
	if err = json.Unmarshal(data, &intVal); err == nil {
		*cts = ControlTaskStatus(intVal)
		return
	}
	if err = json.Unmarshal(data, &strVal); err == nil {
		*cts = taskStatuses[strVal]
		return
	}
	return
}

type taskID struct {
	Url          string            `json:"url"`
	Status       ControlTaskStatus `json:"status"`
	Error        error
	ErrorMessage string `json:"last_error"`
}

type AsyncRequest struct {
	Async bool `json:"async,omitempty"`
}

func (rs *Session) waitAllTasks(tasks []*taskID) error {
	return retry.Do(Timeout, func() error {
		var (
			newTasks []*taskID
			tempErr  error
		)
		watcher := make(chan *taskID)

		for _, task := range tasks {
			task := task
			go func() {
				var ct controlTask
				taskUrl, _ := url.Parse(task.Url)
				err := rs.Request(MethodGet, taskUrl.Path, nil, &ct)
				task.Error = err
				task.Status = ct.Status
				if err != nil {
					task.ErrorMessage = fmt.Sprintf(
						"Couldn't connect to %s",
						task.Url,
					)
				} else {
					task.ErrorMessage = ct.LastError
				}
				watcher <- task
			}()
		}

		for range tasks {
			tid := <-watcher
			if tid.Error != nil {
				// Change this to ControlTaskStatusIncomplete if you
				// want to retry on HTTP error rather than to give up.
				tid.Status = ControlTaskStatusError
			}
			switch tid.Status {
			case ControlTaskStatusSuccess:
				// Task succeeded, we're done with it.
			case ControlTaskStatusError, ControlTaskStatusCanceled:
				return fmt.Errorf("Task <%s> failed due to %v. Cause: %s",
					tid.Url,
					tid.Status,
					tid.ErrorMessage,
				)
			default:
				tempErr = &retry.TemporaryError{
					Err: fmt.Errorf(
						"Task <%s> didn't complete yet: %v",
						tid.Url,
						tid.Status,
					),
				}
				newTasks = append(newTasks, tid)
			}
		}
		tasks = newTasks
		return tempErr
	})
}

func (rs *Session) AsyncRequest(method HttpMethod, uri string, body interface{}) error {
	var tIDs []*taskID
	if body == nil {
		body = &AsyncRequest{Async: true}
	}

	err := rs.Request(method, uri, body, &tIDs)
	if err != nil {
		return err
	}

	if err = rs.waitAllTasks(tIDs); err != nil {
		return err
	}

	return nil
}

type Skipper interface {
	SkipSecurityPrefix()
}

func (rs *Session) Request(method HttpMethod, uri string, body interface{}, result interface{}) error {
	var jsonBody []byte
	var err error

	Log.Debug("request",
		"method", method,
		"uri", uri,
		"body", body,
		"result", result,
	)

	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			Log.Error("http request failed during json.Marshal",
				"err", err, "method", method, "uri", uri, "body", body)
			return errors.New(err)
		}
	}

	resp, resBody, err := rs.requestHttp(method, uri, jsonBody)
	if err != nil {
		Log.Error("http request failed",
			"err", err, "method", method, "uri", uri, "body", body)

		if e, ok := err.(*restError); ok {
			if e.Response.StatusCode == http.StatusUnauthorized {
				Log.Warn("Got eManage Auth Error, relogin and retrying ...")
				rs.init() // invalidate connection before we retry
				err := rs.Login(rs.credentials.User, rs.credentials.Password)
				if err != nil {
					return err
				}
				resp, resBody, err = rs.requestHttp(method, uri, jsonBody)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			return err
		}
	}

	if result == nil {
		if resp.StatusCode != http.StatusNoContent {
			Log.Debug("Request returns content, but no result structure was provided", "method", method, "uri", uri)
		}
		return nil
	}

	switch result.(type) {
	case Skipper:
		Log.Debug("Skipping Security Prefix check")
	default:
		if !bytes.HasPrefix(resBody, securityPrefix) {
			errMsg := fmt.Sprintf("The returned content does not start with the security prefix: %q", securityPrefix)
			return NewRestError(errMsg, resp, resBody)
		}
	}

	resBody = bytes.TrimPrefix(resBody, securityPrefix)

	err = json.Unmarshal(resBody, result)
	if err != nil {
		Log.Error("http request failed during json.Unmarshal",
			"err", err, "method", method, "uri", uri, "body", body, "result", result, "resBody", string(resBody))
		return ejson.NewError(err, resBody)
	}

	Log.Debug("response",
		"method", method,
		"uri", uri,
		"body", string(jsonBody),
		"resp", resp,
		"resBody", string(resBody),
	)

	return nil
}
