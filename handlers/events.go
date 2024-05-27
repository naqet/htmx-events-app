package handlers

import (
	"errors"
	"fmt"
	"htmx-events-app/db"
	"htmx-events-app/internal/chttp"
	"htmx-events-app/middlewares"
	"htmx-events-app/utils"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type eventsHandler struct {
	db *gorm.DB
}

func NewEventsHandler(app *chttp.App) {
	route := app.Group("/events")
	h := eventsHandler{app.DB}

	route.Use(middlewares.Auth)
	route.Get("/{id}", h.getById)
	route.Post("/{$}", h.createEvent)
}

func (h *eventsHandler) getById(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	var event db.Event
	err := h.db.Preload("Hosts", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("Email")
	}).Where("id = ?", id).First(&event).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return chttp.NotFoundError("Event with such ID doesn't exist")
	} else if err != nil {
		return err
	}

	err = utils.WriteJson(w, event)
	return err
}

func (h *eventsHandler) createEvent(w http.ResponseWriter, r *http.Request) error {
	type request struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Place       string    `json:"place"`
		StartDate   time.Time `json:"startDate"`
		EndDate     time.Time `json:"endDate"`
		Hosts       []string  `json:"hosts"`
	}

	var data request

	err := utils.GetDataFromBody(r.Body, &data)

	if err != nil {
		return chttp.BadRequestError()
	}

	email, ok := r.Context().Value("email").(string)

	if !ok || email == "" {
		return fmt.Errorf("Senders email couldn't be obtained from the request")
	}

	err = h.db.Where("title = ?", data.Title).First(&db.Event{}).Error

	if err == nil {
		return chttp.BadRequestError("Event with this title already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	var hosts []*db.User
	for _, email := range data.Hosts {
		hosts = append(hosts, &db.User{Email: email})
	}

	hosts = append(hosts, &db.User{Email: email})

	event := db.Event{
		Title:       data.Title,
		Description: data.Description,
		Place:       data.Place,
		StartDate:   data.StartDate,
		EndDate:     data.EndDate,
		Hosts:       hosts,
	}

	err = h.db.Create(&event).Error

	if err != nil {
		return err
	}

	w.Write([]byte(event.ID))

	return nil
}
