package config

import (
	"fmt"
	"net/http"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// allowedOrigins := []string{"https://mytodolist.my.id", "http://127.0.0.1:5500", "http://127.0.0.1:5501"}
			// origin := r.Header.Get("Origin")

			// for _, allowedOrigin := range allowedOrigins {
			// 	if allowedOrigin == origin {
			// 		w.Header().Set("Access-Control-Allow-Origin", origin)
			// 		break
			// 	}
			// }
			// Set CORS headers for the preflight request
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,Token")
				w.Header().Set("Access-Control-Max-Age", "3600")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Received request:", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
}
