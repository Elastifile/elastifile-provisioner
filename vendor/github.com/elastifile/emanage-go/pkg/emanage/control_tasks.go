package emanage

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/elastifile/emanage-go/pkg/rest"
	"github.com/elastifile/emanage-go/pkg/retry"

	"github.com/elastifile/errors"
)

const (
	controlTasksUri = "api/control_tasks"
)

type controlTasks struct {
	conn *rest.Session
}

type ControlTask struct {
	Attempts     int         `json:"attempts"`
	CreatedAt    string      `json:"created_at"`
	CurrentStep  interface{} `json:"current_step"`
	Host         Host        `json:"host"`
	ID           int         `json:"id"`
	LastError    string      `json:"last_error"`
	Name         string      `json:"name"`
	Priority     int         `json:"priority"`
	Queue        interface{} `json:"queue"`
	Status       string      `json:"status"`
	StepProgress interface{} `json:"step_progress"`
	StepTotal    interface{} `json:"step_total"`
	UpdatedAt    string      `json:"updated_at"`
	UUID         string      `json:"uuid"`
}

func (ct *controlTasks) GetAll(opts *GetAllOpts) (result []ControlTask, err error) {
	err = ct.conn.Request(rest.MethodGet, controlTasksUri, opts, &result)
	return
}

func (ct *controlTasks) GetRecent() (result ControlTask, err error) {
	var r []ControlTask
	r, err = ct.GetRecentSince(nil)
	if len(r) > 0 {
		result = r[0]
	}
	return
}

func (ct *controlTasks) GetRecentSince(opts *RecentOpts) (result []ControlTask, err error) {
	err = ct.conn.Request(rest.MethodGet, filepath.Join(controlTasksUri, "recent"), opts, &result)
	return
}

// Monitors (prints) emanage new control tasks and returns a channel as an handle to be signaled finish
func (ct *controlTasks) Monitor() (chan<- bool, error) {
	pollInterval := 5 * time.Second
	timeout := 1 * time.Minute
	done := make(chan bool)

	recentTask, err := ct.GetRecent()
	if err != nil {
		return done, errors.Errorf("Monitor Control Tasks: cannot start monitoring, err: %v", err)
	}
	go func() {
		for {
			select {
			case <-done:
				break
			default:
				var tasks []ControlTask
				err := retry.Do(timeout, func() (err error) {
					tasks, err = ct.GetRecentSince(&RecentOpts{recentTask.ID})
					if err != nil {
						return &retry.TemporaryError{Err: err}
					}
					return nil
				})
				if err != nil {
					Log.Error("Monitor Control Tasks: failed after some retries", "err", err)
					done <- true
				}

				if len(tasks) > 0 {
					recentTask = tasks[len(tasks)-1] // save last task
					for _, t := range tasks {        // print new tasks
						fmt.Printf("(mgmt control task ID=%v) %v\n", t.ID, t.Name)
					}
				}
			}
			time.Sleep(pollInterval)
		}
	}()

	return done, nil
}
