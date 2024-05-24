package handlers

import (
	"errors"
	"htmx-events-app/db"
	"htmx-events-app/internal/chttp"
	"htmx-events-app/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type authHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) http.Handler {
	route := chttp.New()
	auth := authHandler{db}

	route.HandleFunc("POST /login", auth.login)
	route.HandleFunc("POST /signup", auth.signUp)
	route.HandleFunc("POST /logout", auth.login)

	return route
}

func (h *authHandler) login(w http.ResponseWriter, r *http.Request) error {
    type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
    }

	var data request
    err := utils.GetDataFromBody(r.Body, &data)

    if err != nil {
        return err
    }

    var user db.User
    err = h.db.Where("email = ?", data.Email).First(&user).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return chttp.BadRequestError(chttp.INCORRECT_CREDENTIALS)
    } else if err != nil {
        return err
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))

    if err != nil {
        return chttp.BadRequestError(chttp.INCORRECT_CREDENTIALS)
    }

    //TODO: implement json tokens

	return nil
}

func (h *authHandler) signUp(w http.ResponseWriter, r *http.Request) error {
	type request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var data request
    err := utils.GetDataFromBody(r.Body, &data)

    if err != nil {
        return err
    }

	err = h.db.Where("email = ?", data.Email).First(&db.User{}).Error

	if err == nil {
		return chttp.NewError(http.StatusBadRequest, "User with such email already exists")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user := db.User{Name: data.Name, Email: data.Email, Password: string(pass)}

    err = h.db.Create(&user).Error

    if err != nil {
        return err
    }

    w.Write([]byte(http.StatusText(http.StatusOK)))
	return nil
}
