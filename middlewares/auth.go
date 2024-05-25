package middlewares

import (
	"htmx-events-app/internal/chttp"
	"net/http"
	"strings"
)

func Auth(next http.Handler, path string) chttp.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if strings.HasPrefix(r.URL.Path, path) {
			cookie, err := r.Cookie("Authorization")

			if err != nil || cookie == nil {
				return chttp.BadRequestError()
			}


            //TODO implement json token verification
		}
		next.ServeHTTP(w, r)
		return nil
	}
}
