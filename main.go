package main

import (
	"errors"
	"fmt"
	"htmx-events-app/db"
	"htmx-events-app/handlers"
	"htmx-events-app/internal/chttp"
	"htmx-events-app/middlewares"
	"net/http"
)

func main() {
	database := db.Init()

	app := chttp.New()

	app.Get("/{$}", func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("Hello World"))
		return nil
	})

	app.Handle("/health", handlers.NewHealthHandler())
	app.Handle("/auth", handlers.NewAuthHandler(database))

	app.Use(middlewares.Logger)

	err := app.Listen("localhost:3000")

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
    fmt.Println("Server stopped")
}
