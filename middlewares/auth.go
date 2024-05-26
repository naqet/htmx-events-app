package middlewares

import (
	"context"
	"errors"
	"htmx-events-app/internal/chttp"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.Handler, path string) chttp.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if strings.HasPrefix(r.URL.Path, path) {
			cookie, err := r.Cookie("Authorization")

			if err != nil || cookie == nil {
				return chttp.UnauthorizedError()
			}

			tokenString := cookie.Value

			secret := os.Getenv("SECRET")
			if secret == "" {
				return errors.New("JWT secret is not set")
			}

			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, chttp.UnauthorizedError("Invalid JWT token")
				}
				return []byte(secret), nil
			})

			if err != nil {
				return err
			}

			if !token.Valid {
				return chttp.UnauthorizedError()
			}

			claims, ok := token.Claims.(jwt.MapClaims)

			if !ok {
				return chttp.UnauthorizedError("Invalid JWT token")
			}

            email, err := claims.GetSubject()

			if err != nil {
				return chttp.UnauthorizedError("Invalid JWT token")
			}

			ctx := context.WithValue(r.Context(), "email", email)

			next.ServeHTTP(w, r.WithContext(ctx))
            return nil
		}
		next.ServeHTTP(w, r)
		return nil
	}
}
