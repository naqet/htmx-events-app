package handlers

import (
	"errors"
	"htmx-events-app/db"
	"htmx-events-app/internal/chttp"
	"htmx-events-app/internal/toast"
	"htmx-events-app/middlewares"
	"htmx-events-app/utils"
	vcomponents "htmx-events-app/views/components"
	vevents "htmx-events-app/views/events"
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
	route.Get("/", h.homePage)
	route.Get("/create", h.createEventPage)
	route.Get("/{title}", h.getByTitle)
	route.Post("/{$}", h.createEvent)
	route.Post("/{title}/agenda-point", h.createAgendaPoint)
}

func (h *eventsHandler) createEventPage(w http.ResponseWriter, r *http.Request) error {
	return vevents.CreateEventPage().Render(r.Context(), w)
}

func (h *eventsHandler) homePage(w http.ResponseWriter, r *http.Request) error {
	email, err := utils.GetEmailFromContext(r)

	if err != nil {
		return err
	}

	var events []db.Event
	err = h.db.
		Joins("JOIN attended_events ON attended_events.event_id = events.id").
		Joins("JOIN users ON users.email = attended_events.user_email").
		Where("users.email = ?", email).
		Preload("Hosts").
		Find(&events).
		Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return vevents.Page(events).Render(r.Context(), w)
}

func (h *eventsHandler) getByTitle(w http.ResponseWriter, r *http.Request) error {
	title := r.PathValue("title")

	var event db.Event
	err := h.db.Where("title = ?", title).Preload("Hosts").Preload("Agenda", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("start_time ASC")
	}).Preload("Attendees").Find(&event).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return chttp.NotFoundError("Event with such title doesn't exist")
	} else if err != nil {
		return err
	}

	email, err := utils.GetEmailFromContext(r)

	if err != nil && !errors.Is(err, utils.ErrEmptyEmail) {
		return err
	}

	var isOwner bool

	for _, host := range event.Hosts {
		if host.Email == email {
			isOwner = true
			break
		}
	}

	return vevents.Details(event, isOwner).Render(r.Context(), w)
}

func (h *eventsHandler) createEvent(w http.ResponseWriter, r *http.Request) error {
	type request struct {
		Title              string          `json:"title"`
		Description        string          `json:"description"`
		Place              string          `json:"place"`
		StartDate          utils.Time      `json:"startDate"`
		EndDate            utils.Time      `json:"endDate"`
		Hosts              utils.StringArr `json:"hosts"`
		AgendaTitles       utils.StringArr `json:"agendaTitles"`
		AgendaDates        utils.TimeArr   `json:"agendaDates"`
		AgendaDescriptions utils.StringArr `json:"agendaDescriptions"`
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

	points, err := handleAgendaPoints(data.AgendaTitles, data.AgendaDates, data.AgendaDescriptions)

	hosts = append(hosts, &db.User{Email: email})

	event := db.Event{
		Title:       data.Title,
		Description: data.Description,
		Place:       data.Place,
		StartDate:   time.Time(data.StartDate),
		EndDate:     time.Time(data.EndDate),
		Hosts:       hosts,
		Attendees:   hosts,
		Agenda:      points,
	}

	err = h.db.Create(&event).Error

	if err != nil {
		return err
	}

	if utils.IsHtmxRequest(r) {
		toast.AddToast(w, toast.SUCCESS, "Event has been added")
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(event.ID))

	return nil
}

func handleAgendaPoints(titles []string, dates []utils.Time, descriptions []string) ([]db.AgendaPoint, error) {
	var result []db.AgendaPoint

	if len(titles) != len(dates) || len(titles) != len(descriptions) || len(dates) != len(descriptions) {
		return result, errors.New("Invalid agenda points entries")
	}

	for i := range len(titles) {
		point := db.AgendaPoint{
			Base:        db.Base{},
			Title:       titles[i],
			Description: descriptions[i],
			StartTime:   time.Time(dates[i]),
		}

		result = append(result, point)
	}

	return result, nil
}

func (h *eventsHandler) createAgendaPoint(w http.ResponseWriter, r *http.Request) error {
	type request struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		StartTime   utils.Time `json:"startTime"`
	}

	var data request

	err := utils.GetDataFromBody(r.Body, &data)

	if err != nil {
		return chttp.BadRequestError()
	}

	title := r.PathValue("title")

	var event db.Event
	err = h.db.Preload("Hosts").Where("title = ?", title).First(&event).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return chttp.NotFoundError("Event with such title doesn't exist")
	} else if err != nil {
		return err
	}

	inTimeSpan := utils.InTimeSpan(event.StartDate, event.EndDate, time.Time(data.StartTime))

	if !inTimeSpan {
		msg := "Agenda point must be in the time span of the event"
		err = toast.AddToast(w, toast.DANGER, msg)
		if err != nil {
			return err
		}
		return chttp.BadRequestError(msg)
	}

	email, err := utils.GetEmailFromContext(r)

	if err != nil {
		return err
	}

	var isOwner bool

	for _, host := range event.Hosts {
		if host.Email == email {
			isOwner = true
			break
		}
	}

	if !isOwner {
		return chttp.UnauthorizedError("Only event's host can add agenda points")
	}

	agendaPoint := db.AgendaPoint{
		Title:       data.Title,
		Description: data.Description,
		StartTime:   time.Time(data.StartTime),
	}

	err = h.db.Model(&event).Association("Agenda").Append(&agendaPoint)

	if err != nil {
		return err
	}

	err = h.db.Preload("Agenda", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("start_time ASC")
	}).Where("title = ?", title).First(&event).Error

	if err != nil {
		return err
	}

	return vcomponents.AgendaList(utils.OrganizeAgendaPoints(event.Agenda), isOwner).Render(r.Context(), w)
}
