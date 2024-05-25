package chttp

import (
	"log/slog"
	"net/http"
)

const (
    INCORRECT_CREDENTIALS = "Incorrect email or password"
)

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

func BadRequestError(args ...string) StatusError {
    msg := http.StatusText(http.StatusBadRequest)
    if len(args) == 1 {
        msg = args[0]
    }
	return StatusError{http.StatusBadRequest, msg}
}

func UnauthorizedError(args ...string) StatusError {
    msg := http.StatusText(http.StatusUnauthorized)
    if len(args) == 1 {
        msg = args[0]
    }
	return StatusError{http.StatusUnauthorized, msg}
}

func withErrorHandling(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		if err != nil {
            slog.Error(err.Error())
			switch e := err.(type) {
			case HttpError:
				http.Error(w, e.Error(), e.Status())
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
	}
}

