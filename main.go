package main

import (
	"errors"
	"fmt"
	"htmx-events-app/db"
	"htmx-events-app/handlers"
	"htmx-events-app/internal/cenv"
	"htmx-events-app/internal/chttp"
	"htmx-events-app/internal/toast"
	"htmx-events-app/middlewares"
	"net/http"
)

func main() {
	cenv.Init()
	database := db.Init()
	app := chttp.New(database)

	app.Use(middlewares.Logger)

	app.Handle("/static", http.FileServer(http.Dir("./static")))

    app.Post("/testing", func(w http.ResponseWriter, r *http.Request) error {
        toast.AddToast(w, toast.SUCCESS, "Testing")
        w.Write([]byte("OK"))
        return nil
    })

	handlers.NewWelcomeHandler(app)
	handlers.NewHealthHandler(app)
	handlers.NewAuthHandler(app)
	handlers.NewDashboardHandler(app)
	handlers.NewEventsHandler(app)
	handlers.NewComponentsHandler(app)
	handlers.NewInvitationsHandler(app)
	handlers.NewAgendaPointsHandler(app)

	fmt.Println("Starting server...")
	err := app.Listen("localhost:3000")
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
	fmt.Println("Server stopped")
}
