package handlers

import (
	"errors"
	"htmx-events-app/db"
	"htmx-events-app/internal/chttp"
	"htmx-events-app/internal/toast"
	"htmx-events-app/middlewares"
	"htmx-events-app/utils"
	"net/http"

	"gorm.io/gorm"
)

type invitationsHandler struct {
	db *gorm.DB
}

func NewInvitationsHandler(app *chttp.App) {
	route := app.Group("/invitations")
	h := invitationsHandler{app.DB}

	route.Use(middlewares.Auth)

	route.Post("/", h.create)
	route.Post("/{id}/accept", h.accept)
	route.Post("/{id}/reject", h.reject)
}

// TODO: implement logic
func (h *invitationsHandler) accept(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *invitationsHandler) reject(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *invitationsHandler) create(w http.ResponseWriter, r *http.Request) error {
	type request struct {
		To      string `json:"to"`
		Event   string `json:"event"`
		Message string `json:"message"`
	}

	var data request
	err := utils.GetDataFromBody(r.Body, &data)

	if err != nil {
		return err
	}

	email, err := utils.GetEmailFromContext(r)

	if err != nil {
		return err
	}

	var receiver db.User
	err = h.db.Where("email = ?", data.To).First(&receiver).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return chttp.BadRequestError("Receiver does not exist")
	}

	if err != nil {
		return err
	}

	var event db.Event
	err = h.db.Preload("Hosts").Where("title = ?", data.Event).First(&event).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return chttp.BadRequestError("Event with such title does not exist")
	}

	if err != nil {
		return err
	}

	var isHost bool
	for _, host := range event.Hosts {
		if host.Email == email {
			isHost = true
			break
		}
	}

	if !isHost {
		return chttp.UnauthorizedError("Only event's hosts can send invitations")
	}

	invitation := db.Invitation{
		From:    db.User{Email: email},
		To:      receiver,
		Event:   event,
		Message: data.Message,
	}

	err = h.db.Create(&invitation).Error

	if err != nil {
		return err
	}

	toast.AddToast(w, toast.SUCCESS, "Invitation has been sent")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(http.StatusText(http.StatusCreated)))
	return nil
}
