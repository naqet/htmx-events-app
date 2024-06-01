package handlers

import (
	"htmx-events-app/internal/chttp"
	vwelcome "htmx-events-app/views/welcome"
	"net/http"
)

type welcomeHandler struct{}

func NewWelcomeHandler(app *chttp.App) {
	h := welcomeHandler{}

	app.Get("/{$}", h.homePage)
}

func (h *welcomeHandler) homePage(w http.ResponseWriter, r *http.Request) error {
	return vwelcome.Home().Render(r.Context(), w)
}
