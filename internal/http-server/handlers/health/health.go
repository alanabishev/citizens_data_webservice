// Package health provides HTTP handlers for health checks.
package health

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

// HealthResponse is the response structure for the health check endpoint.
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

// Check is an HTTP handler function for health checks.
// It returns the service status and basic information.
// This endpoint can be used by load balancers and monitoring systems.
func Check(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("health check requested")

		render.JSON(w, r, HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Service:   "citizens-data-webservice",
			Version:   "1.0.0",
		})
	}
}
