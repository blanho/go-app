// internal/middleware/logging.go
package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)
		
		duration := time.Since(start)
		
		log.Printf(
			"%s %s %d %s",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration.String(),
		)
	})
}