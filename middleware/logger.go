package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger is a middleware that logs the incoming request.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the request details
		log.Printf("%s /%s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
