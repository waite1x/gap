package common

import (
	"fmt"
)

type HttpResult struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

const ErrCodeUnkown string = "UnkonwError"

type HttpError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewHttpError(err error) HttpError {
	return HttpError{
		Message: err.Error(),
		Code:    ErrCodeUnkown,
	}
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("http error: code: %s, msg: %s, httpstatus:%d", e.Code, e.Message, e.Code)
}
