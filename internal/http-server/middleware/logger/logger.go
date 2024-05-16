// Package logger provides a middleware for logging HTTP requests.
package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

// New is a function that creates a new logging middleware.
// It takes a logger as a parameter and returns a middleware function.
// The middleware function logs the method, path, remote address, user agent, request ID, status, bytes written, and duration of each request.
func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Add component information to the logger
		log := log.With(
			slog.String("component", "middleware/logger"),
		)

		// Log that the logger middleware is enabled
		log.Info("logger middleware enabled")

		// Define the middleware function
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Create a new log entry with request information
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			// Create a new response writer that allows us to get the status and bytes written
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Record the start time of the request
			t1 := time.Now()
			// When the request is done, log the status, bytes written, and duration
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			// Call the next middleware or handler
			next.ServeHTTP(ww, r)
		}

		// Return the middleware function
		return http.HandlerFunc(fn)
	}
}
