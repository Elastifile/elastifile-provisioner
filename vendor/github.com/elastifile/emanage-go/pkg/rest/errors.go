package rest

import (
	"fmt"
	"net/http"
)

type restError struct {
	Description string
	Response    *http.Response
	Body        string
}

func NewRestError(description string, response *http.Response, resBody []byte) *restError {
	err := &restError{
		Description: description,
		Response:    response,
		Body:        string(resBody),
	}
	return err
}

func (e *restError) Error() string {
	return fmt.Sprint(e.Description,
		"\n\nHTTP response:\n", e.Response,
		"\n\nBody:\n", e.Body,
	)
}
