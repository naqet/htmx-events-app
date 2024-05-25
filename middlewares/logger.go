package middlewares

import (
	"htmx-events-app/internal/chttp"
	"log"
	"net/http"
	"strings"
)

func Logger(next http.Handler, path string) chttp.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if !strings.HasPrefix(r.URL.Path, "/static") && strings.HasPrefix(r.URL.Path, path) {
			log.Printf("[%s] %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		}

		next.ServeHTTP(w, r)
		return nil
	}
}
