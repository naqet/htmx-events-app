package handlers

import (
	"htmx-events-app/internal/chttp"
	"net/http"
)

type healthHandler struct{}

func NewHealthHandler(app *chttp.App) {
	route := app.Group("/health")
	h := healthHandler{}

	route.Get("/", h.healthy)
}

func (h *healthHandler) healthy(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte(http.StatusText(http.StatusOK)))
	return nil
}
