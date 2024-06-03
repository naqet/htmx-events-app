package handlers

import (
	"htmx-events-app/db"
	"htmx-events-app/internal/chttp"
	"htmx-events-app/middlewares"
	"htmx-events-app/utils"
	vdashboard "htmx-events-app/views/dashboard"
	"net/http"

	"gorm.io/gorm"
)

type dashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(app *chttp.App) {
	route := app.Group("/dashboard")
	h := dashboardHandler{app.DB}

	route.Use(middlewares.Auth)

	route.Get("/{$}", h.homePage)
}

func (h *dashboardHandler) homePage(w http.ResponseWriter, r *http.Request) error {
	email, err := utils.GetEmailFromContext(r)

	if err != nil {
		return err
	}

	var invitations []db.Invitation
	limit := 5
	var count int64
	err = h.db.
		Preload("From").
		Preload("Event").
		Where("to_email = ?", email).
		Limit(limit).
		Find(&invitations).
		Count(&count).
		Error

	if err != nil {
		return err
	}

	return vdashboard.Page(invitations, int(count) > limit).Render(r.Context(), w)
}
