package handlers

import (
	"errors"
	"htmx-events-app/db"
	"htmx-events-app/internal/chttp"
	"htmx-events-app/middlewares"
	"htmx-events-app/utils"
	"net/http"

	"gorm.io/gorm"
)

type agendaPointsHandler struct {
	db *gorm.DB
}

func NewAgendaPointsHandler(app *chttp.App) {
	route := app.Group("/agenda-points")
	h := agendaPointsHandler{app.DB}

	route.Use(middlewares.Auth)

	route.Delete("/{id}", h.delete)
}

func (h *agendaPointsHandler) delete(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	email, err := utils.GetEmailFromContext(r)

	if err != nil {
		return err
	}

	point := db.AgendaPoint{}
	err = h.db.Preload("Event.Hosts").Where("id = ?", id).First(&point).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return chttp.BadRequestError("Agenda point with such ID does't exist")
	} else if err != nil {
		return err
	}

	var isOwner = utils.IsEventOwner(email, point.Event)

	if !isOwner {
		return chttp.UnauthorizedError("Only event's owner can change the agenda")
	}

	err = h.db.Delete(&point).Error

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
