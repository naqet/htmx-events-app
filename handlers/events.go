package handlers

import (
	"errors"
	"htmx-events-app/db"
	"htmx-events-app/internal/chttp"
	"htmx-events-app/internal/toast"
	"htmx-events-app/middlewares"
	"htmx-events-app/utils"
	vevents "htmx-events-app/views/events"
	"net/http"

	"gorm.io/gorm"
)

type eventsHandler struct {
	db *gorm.DB
}

func NewEventsHandler(app *chttp.App) {
	route := app.Group("/events")
	h := eventsHandler{app.DB}

	route.Use(middlewares.Auth)
	route.Get("/", h.homePage)
	route.Get("/{title}", h.getById)
	route.Post("/{$}", h.createEvent)
}

func (h *eventsHandler) homePage(w http.ResponseWriter, r *http.Request) error {
	email, err := utils.GetEmailFromContext(r)

	if err != nil {
		return err
	}

	var events []db.Event
	err = h.db.
		Joins("JOIN hosted_events ON hosted_events.event_id = events.id").
		Joins("JOIN users ON users.email = hosted_events.user_email").
		Where("users.email = ?", email).
		Preload("Hosts").Find(&events).
		Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return vevents.Page(events).Render(r.Context(), w)
}

func (h *eventsHandler) getById(w http.ResponseWriter, r *http.Request) error {
	title := r.PathValue("title")

	var event db.Event
	err := h.db.Preload("Hosts").Where("title = ?", title).First(&event).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return chttp.NotFoundError("Event with such ID doesn't exist")
	} else if err != nil {
		return err
	}

	return vevents.Details(event).Render(r.Context(), w)
}

func (h *eventsHandler) createEvent(w http.ResponseWriter, r *http.Request) error {
	type request struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Place       string     `json:"place"`
		StartDate   utils.Time `json:"startDate"`
		EndDate     utils.Time `json:"endDate"`
		Hosts       []string   `json:"hosts"`
	}

	var data request

	err := utils.GetDataFromBody(r.Body, &data)

	if err != nil {
		return chttp.BadRequestError()
	}

	email, err := utils.GetEmailFromContext(r)

	if err != nil {
		return err
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
		StartDate:   data.StartDate.Time,
		EndDate:     data.EndDate.Time,
		Hosts:       hosts,
	}

	err = h.db.Create(&event).Error

	if err != nil {
		return err
	}

	if utils.IsHtmxRequest(r) {
		toast.AddToast(w, toast.SUCCESS, "Event has been added")
	}

	w.Write([]byte(event.ID))

	return nil
}
