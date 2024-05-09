package config

import "net/http"

// CorsMiddleware dengan metode dan origin yang diizinkan.
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Tetapkan origin yang diizinkan
			allowedOrigins := map[string]bool{
				"https://mytodolist.my.id": true,
				"http://127.0.0.1:5500":    true,
				"http://127.0.0.1:5501":    true,
			}

			origin := r.Header.Get("Origin")

			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Max-Age", "3600")

				// Jika preflight
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusNoContent)
					return
				}
			} else {
				w.WriteHeader(http.StatusForbidden) // Origin tidak diizinkan
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}
