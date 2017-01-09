package ejson

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Error struct {
	Err     error
	RawData []byte
	ErrInfo string
}

func NewError(err error, rawData []byte) *Error {
	var errInfo string
	if err == strconv.ErrSyntax {
		err = fmt.Errorf("strconv: %v", err)
	} else if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
		err = fmt.Errorf("json: offset: %v: %v", jsonErr.Offset, err)
		errInfo = string(rawData)[:jsonErr.Offset] + "<=="
	} else if jsonErr, ok := err.(*json.SyntaxError); ok {
		err = fmt.Errorf("json: offset: %v: %v", jsonErr.Offset, err)
		errInfo = string(rawData)[:jsonErr.Offset] + "<=="
	}
	return &Error{err, rawData, errInfo}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%T: %T: %s\nRawData: %s\nErrInfo: %s", e, e.Err, e.Err.Error(), e.RawData, e.ErrInfo)
}
