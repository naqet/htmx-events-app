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
	route.Post("/{$}", h.createEvent)
}

func (h *eventsHandler) createEvent(w http.ResponseWriter, r *http.Request) error {
	type request struct {
		Title            string    `json:"title"`
		Description      string    `json:"description"`
		Place            string    `json:"place"`
		StartDate        time.Time `json:"startDate"`
		EndDate          time.Time `json:"endDate"`
		AdditionalOwners []string  `json:"additionalOwners"`
	}

	var data request

	err := utils.GetDataFromBody(r.Body, &data)

	if err != nil {
		return chttp.BadRequestError()
	}

	senderUserId, ok := r.Context().Value("id").(string)

	if !ok || senderUserId == "" {
		return fmt.Errorf("Senders userID couldn't be obtained from the request")
	}

    var user db.User

    err = h.db.Where("id = ?", senderUserId).First(&user).Error

    if err != nil {
        return err
    }

    err = h.db.Where("title = ?", data.Title).First(&db.Event{}).Error

    if err == nil {
        return chttp.BadRequestError("Event with this title already exists")
    }

    if !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }

    var owners []*db.User

    for _, email := range data.AdditionalOwners {
        owners = append(owners, &db.User{Email: email})
    }

    event := db.Event{
    	Title:       data.Title,
    	Description: data.Description,
    	Place:       data.Place,
    	StartDate:   data.StartDate,
    	EndDate:     data.EndDate,
    	Owners:      owners,
    }

    err = h.db.Create(&event).Error

	return err
}
