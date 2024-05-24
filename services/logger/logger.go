package logger

import (
	"log"
	"net/http"
)

type logger struct{}

func NewLogger() *logger {
	return &logger{}
}

func (l *logger) Middleware(next http.HandlerFunc) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("[%s] %s %s", r.RemoteAddr, r.Method, r.URL.Path)
        next.ServeHTTP(w,r)
    })
}
