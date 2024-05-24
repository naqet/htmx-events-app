package chttp

import (
	"fmt"
	"net/http"
)

type HandlerFunc = func(w http.ResponseWriter, r *http.Request) error
type Middleware = func(next http.Handler) http.Handler

type app struct {
	path        string
	mux         *http.ServeMux
	middlewares []Middleware
}

func New() *app {
    return &app{path: "", mux: http.NewServeMux()}
}

func (a app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *app) Handle(path string, handler http.Handler) {
	a.mux.Handle(a.path+path+"/", http.StripPrefix(a.path+path, handler))
}

func (a *app) Group(path string) *app {
	return &app{
		path:        a.path + path,
		mux:         a.mux,
		middlewares: a.middlewares,
	}
}

func (a *app) Use(middle Middleware) {
	a.middlewares = append(a.middlewares, middle)
}

func (a *app) Get(path string, f HandlerFunc) {
    pattern := fmt.Sprintf("GET %s", a.path + path)
	a.mux.HandleFunc(pattern, withErrorHandling(f))
}

func (a *app) Post(path string, f HandlerFunc) {
    pattern := fmt.Sprintf("POST %s", a.path + path)
	a.mux.HandleFunc(pattern, withErrorHandling(f))
}

func (a *app) Put(path string, f HandlerFunc) {
    pattern := fmt.Sprintf("PUT %s", a.path + path)
	a.mux.HandleFunc(pattern, withErrorHandling(f))
}

func (a *app) Delete(path string, f HandlerFunc) {
    pattern := fmt.Sprintf("DELETE %s", a.path + path)
	a.mux.HandleFunc(pattern, withErrorHandling(f))
}

func (a *app) Listen(addr string) error {
	var handler http.Handler = a.mux
	for _, middleware := range a.middlewares {
		handler = middleware(handler)
	}
	return http.ListenAndServe(addr, handler)
}
