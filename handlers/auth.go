package handlers

import (
	"errors"
	"htmx-events-app/db"
	"htmx-events-app/internal/chttp"
	"htmx-events-app/utils"
	vauth "htmx-events-app/views/auth"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type authHandler struct {
	db *gorm.DB
}

func NewAuthHandler(app *chttp.App) {
	route := app.Group("/auth")
	auth := authHandler{app.DB}

    route.Get("/login", auth.loginPage)
    route.Get("/signup", auth.signUpPage)

	route.Post("/login", auth.login)
	route.Post("/signup", auth.signUp)
	route.Post("/logout", auth.login)
}

func (h *authHandler) loginPage(w http.ResponseWriter, r *http.Request) error {
    return vauth.LoginPage().Render(r.Context(), w)
}

func (h *authHandler) signUpPage(w http.ResponseWriter, r *http.Request) error {
    return vauth.SignUpPage().Render(r.Context(), w)
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

	expiration := time.Now().Add(time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.Email,
		"time": expiration.Unix(),
	})

	secret := os.Getenv("SECRET")

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Path:     "/",
		Expires:  expiration,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &cookie)
    w.Write([]byte(http.StatusText(http.StatusOK)))

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
