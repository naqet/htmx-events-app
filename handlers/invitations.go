package handlers

import (
	"htmx-events-app/internal/chttp"
	"net/http"

	"gorm.io/gorm"
)

type invitationsHandler struct {
	db *gorm.DB
}

func NewInvitationsHandler(app *chttp.App) {
	route := app.Group("/invitations")
	h := invitationsHandler{app.DB}

    route.Post("/", h.create)
}

func (h *invitationsHandler) create(w http.ResponseWriter, r *http.Request) error {
    return nil
}
