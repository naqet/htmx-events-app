package handlers

import (
	"errors"
	"htmx-events-app/db"
	"htmx-events-app/internal/chttp"
	"htmx-events-app/middlewares"
	"htmx-events-app/utils"
	vcomponents "htmx-events-app/views/components"
	"net/http"
	"slices"

	"gorm.io/gorm"
)

type componentsHandler struct {
	db *gorm.DB
}

func NewComponentsHandler(app *chttp.App) {
	route := app.Group("/components")
	h := componentsHandler{app.DB}

	route.Use(middlewares.Auth)

	route.Post("/users", h.users)
}

func (h *componentsHandler) users(w http.ResponseWriter, r *http.Request) error {
	email, err := utils.GetEmailFromContext(r)

	if err != nil {
		return err
	}

	search := r.FormValue("search")
    hosts := r.Form["hosts"]

	var users []db.User
	err = h.db.Where("email <> ? AND name LIKE ?", email, "%" + search + "%").Or("email IN ?", hosts).Find(&users).Error

    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }

	var options []vcomponents.Option
	for _, user := range users {
        option := vcomponents.Option{
			Label: user.Name,
			Value: user.Email,
        }

        if slices.Index(hosts, user.Email) != -1 {
            option.Checked = true
        }
		options = append(options, option)
	}
	return vcomponents.MultiselectOptions(options).Render(r.Context(), w)
}
