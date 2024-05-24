package chttp

import (
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

func (a *app) HandleFunc(path string, f HandlerFunc) {
	a.mux.HandleFunc(a.path+path, withErrorHandling(f))
}

func withErrorHandling(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		if err != nil {
			switch e := err.(type) {
			case HttpError:
				http.Error(w, e.Error(), e.Status())
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
	}
}

func (a *app) Listen(addr string) error {
	var handler http.Handler = a.mux
	for _, middleware := range a.middlewares {
		handler = middleware(handler)
	}
	return http.ListenAndServe(addr, handler)
}
