package chttp

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type HandlerFunc = func(w http.ResponseWriter, r *http.Request) error
type Middleware = func(next http.Handler, path string) HandlerFunc

type App struct {
	path        string
	mux         *http.ServeMux
	middlewares *map[string][]Middleware
	DB          *gorm.DB
}

func New(db *gorm.DB) *App {
	return &App{
		path:        "",
		mux:         http.NewServeMux(),
		middlewares: &map[string][]Middleware{},
		DB:          db,
	}
}

func (a App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func (a *App) Handle(path string, handler http.Handler) {
	a.mux.Handle(a.path+path+"/", http.StripPrefix(a.path+path, handler))
}

func (a *App) Group(path string) *App {
	return &App{
		path:        a.path + path,
		mux:         a.mux,
		middlewares: a.middlewares,
	}
}

func (a *App) Use(middle Middleware) {
	(*a.middlewares)[a.path] = append((*a.middlewares)[a.path], middle)
}

func (a *App) Get(path string, f HandlerFunc) {
	if a.path != "/" {
		path = a.path + path
	}

	pattern := fmt.Sprintf("GET %s", path)
	a.mux.HandleFunc(pattern, withErrorHandling(f))
}

func (a *App) Post(path string, f HandlerFunc) {
	if a.path != "/" {
		path = a.path + path
	}

	pattern := fmt.Sprintf("POST %s", path)
	a.mux.HandleFunc(pattern, withErrorHandling(f))
}

func (a *App) Put(path string, f HandlerFunc) {
	if a.path != "/" {
		path = a.path + path
	}

	pattern := fmt.Sprintf("PUT %s", path)
	a.mux.HandleFunc(pattern, withErrorHandling(f))
}

func (a *App) Delete(path string, f HandlerFunc) {
	if a.path != "/" {
		path = a.path + path
	}

	pattern := fmt.Sprintf("DELETE %s", path)
	a.mux.HandleFunc(pattern, withErrorHandling(f))
}

func (a *App) Listen(addr string) error {
	var handler http.Handler = a.mux
	for path, middlewares := range *a.middlewares {
		for _, middleware := range middlewares {
			handler = withErrorHandling(middleware(handler, path))
		}
	}
	return http.ListenAndServe(addr, handler)
}
