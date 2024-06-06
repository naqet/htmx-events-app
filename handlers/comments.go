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

type commentsHandler struct {
	db *gorm.DB
}

func NewCommentsHandler(app *chttp.App) {
	route := app.Group("/comments")
	h := commentsHandler{app.DB}

	route.Use(middlewares.Auth)

	route.Delete("/{id}", h.delete)
}

func (h *commentsHandler) delete(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	email, err := utils.GetEmailFromContext(r)

	if err != nil {
		return err
	}

	comment := db.Comment{}
	err = h.db.Where("id = ?", id).First(&comment).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return chttp.BadRequestError("Comment with such ID does't exist")
	} else if err != nil {
		return err
	}

	if comment.FromEmail != email {
		return chttp.UnauthorizedError("Only comment's author can delete it")
	}

	err = h.db.Delete(&comment).Error

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
