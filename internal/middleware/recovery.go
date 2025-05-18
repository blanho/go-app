// internal/middleware/recovery.go
package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/yourusername/azure-go-app/internal/telemetry"
)

func RecoveryMiddleware(tel *telemetry.Telemetry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					stackTrace := debug.Stack()
					errMsg := fmt.Sprintf("PANIC: %v\n%s", err, stackTrace)
					
					tel.TrackException(fmt.Errorf("%v", err), map[string]string{
						"stack_trace": string(stackTrace),
						"url":         r.URL.String(),
						"method":      r.Method,
					})
					
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			
			next.ServeHTTP(w, r)
		})
	}
}