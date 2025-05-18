// internal/middleware/telemetry.go
package middleware

import (
	"net/http"
	"time"

	"github.com/yourusername/azure-go-app/internal/telemetry"
)

func TelemetryMiddleware(tel *telemetry.Telemetry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			
			rw := newResponseWriter(w)
			
			next.ServeHTTP(rw, r)
			
			duration := time.Since(startTime)
			success := rw.statusCode < 400
			
			tel.TrackRequest(
				r.Context(),
				r.Method+" "+r.URL.Path,
				startTime,
				duration,
				string(rw.statusCode),
				success,
			)
			
			tel.TrackMetric("http.request.duration", float64(duration.Milliseconds()), map[string]string{
				"path":   r.URL.Path,
				"method": r.Method,
				"status": string(rw.statusCode),
			})
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}