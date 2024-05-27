package handlers

import (
	"htmx-events-app/internal/chttp"
	"htmx-events-app/middlewares"
	vdashboard "htmx-events-app/views/dashboard"
	"net/http"
)

type dashboardHandler struct{}

func NewDashboardHandler(app *chttp.App){
    route := app.Group("/dashboard")
    h := dashboardHandler{}

    route.Use(middlewares.Auth)

    route.Get("/{$}", h.homePage)
}

func (h *dashboardHandler) homePage(w http.ResponseWriter, r *http.Request) error {
    return vdashboard.Page().Render(r.Context(), w)
}
