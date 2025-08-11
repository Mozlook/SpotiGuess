package middleware

import (
	"net/http"
	"os"
)

// EnableCORS is a middleware that enables Cross-Origin Resource Sharing (CORS)
// for incoming HTTP requests.
//
// It sets the following headers on all responses:
//
//   - Access-Control-Allow-Origin: http://127.0.0.1:5173
//   - Access-Control-Allow-Methods: GET, POST, OPTIONS
//   - Access-Control-Allow-Headers: Content-Type, Authorization
//
// If the incoming request method is OPTIONS (CORS preflight),
// the middleware responds immediately with HTTP 200 and does not invoke the next handler.
//
// For all other requests, the wrapped handler is invoked normally.
//
// Use this middleware to allow frontend applications running on
// http://127.0.0.1:5173 (e.g. Vite/React) to access your backend API
// during local development.
func EnableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWED_CORS"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
