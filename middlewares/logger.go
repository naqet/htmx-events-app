package middlewares

import (
	"log"
	"net/http"
)

func Logger (next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("[%s] %s %s", r.RemoteAddr, r.Method, r.URL.Path)
        next.ServeHTTP(w,r)
    })
}