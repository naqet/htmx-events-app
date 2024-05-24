package main

import (
	"htmx-events-app/services/logger"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	log := logger.NewLogger()

	mux.Handle("/{$}", log.Middleware(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	}))

	mux.Handle("/health", log.Middleware(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(http.StatusText(http.StatusOK)))
	}))

	mux.Handle("/", log.Middleware(http.NotFound))

	http.ListenAndServe("localhost:3000", mux)
}
