// Package save provides HTTP handlers for saving person information.
package save

import (
	"citizen_webservice/internal/http-server/handlers/request_validator"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

// Request is the structure for the request body of the Person handler.
type Request struct {
	IIN   string `json:"iin" validate:"required,len=12,iin"` // Individual Identification Number
	Name  string `json:"name" validate:"required"`           // Name of the person
	Phone string `json:"phone" validate:"required"`          // Phone number of the person
}

// PersonSaver is an interface for saving person information.
type PersonSaver interface {
	SavePerson(iin string, name string, phone string) error
}

// PersonResponse is the response structure for the Person handler.
type PersonResponse struct {
	Success bool     `json:"success"` // Indicates if the operation was successful
	Errors  []string `json:"errors"`  // List of error messages, if any
}

// Person is a HTTP handler function for saving a person's information.
// It decodes the request body, validates the request, saves the person information,
// and returns a JSON response.
func Person(log *slog.Logger, personSaver PersonSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.Person"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			handleError(w, r, log, err, "Failed to decode request body")
			return
		}

		customValidator := request_validator.GetValidator()
		log.Info("request body decoded", slog.Any("request", req))
		if err := customValidator.Struct(req); err != nil {
			handleError(w, r, log, err, "Validation failed")
			return
		}

		err = personSaver.SavePerson(req.IIN, req.Name, req.Phone)
		if err != nil {
			handleError(w, r, log, err, "Failed to save person")
			return
		}

		log.Info("person added", slog.String("id", req.IIN))
		render.JSON(w, r, PersonResponse{
			Success: true,
		})
	}
}

// handleError is a helper function to handle errors.
// It logs the error, determines the appropriate HTTP status code,
// and sends a JSON response with the error message.
func handleError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error, message string) {
	log.Error(message, Err(err))
	status := http.StatusInternalServerError
	if errors.Is(err, io.EOF) || request_validator.CheckErrorIsValidation(err) {
		status = http.StatusBadRequest
	}
	render.Status(r, status)
	render.JSON(w, r, PersonResponse{
		Success: false,
		Errors:  []string{fmt.Sprintf("%s: %s", message, err.Error())},
	})
}

// Err is a helper function to create a structured log attribute for errors.
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
