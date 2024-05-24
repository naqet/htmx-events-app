package handlers

import (
	"htmx-events-app/internal/chttp"
	"net/http"
)

type healthHandler struct { }

func NewHealthHandler() http.Handler {
	app := chttp.New()
    h := healthHandler{}

	app.Get("/", h.healthy)

	return app
}

func (h *healthHandler) healthy(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte(http.StatusText(http.StatusOK)))
	return nil
}
