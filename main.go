package main

import (
	"htmx-events-app/db"
	"htmx-events-app/handlers"
	"htmx-events-app/internal/chttp"
	"htmx-events-app/middlewares"
	"net/http"
)

func main() {
    database := db.Init()

	app := chttp.New()

	app.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("Hello World"))
		return nil
	})

    app.Handle("/health", handlers.NewHealthHandler())
    app.Handle("/auth", handlers.NewAuthHandler(database))

	app.Use(middlewares.Logger)

	app.Listen("localhost:3000")

    http.ListenAndServe("localhost:3000", app)
}
