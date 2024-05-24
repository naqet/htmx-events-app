package chttp

import "net/http"

type HttpError interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  string
}

func (s StatusError) Error() string {
	return s.Err
}

func (s StatusError) Status() int {
	return s.Code
}

func NewError(code int, msg string) StatusError {
	return StatusError{code, msg}
}

func BadRequestError() StatusError {
	return StatusError{http.StatusBadRequest, http.StatusText(http.StatusBadRequest)}
}
