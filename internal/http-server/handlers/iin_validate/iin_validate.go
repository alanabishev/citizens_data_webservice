// Package iin_validate provides HTTP handlers for validating IIN (Individual Identification Number).
package iin_validate

import (
	resp "citizen_webservice/internal/http-server/handlers/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"

	"citizen_webservice/internal/iin_validator"
	"github.com/go-chi/render"
)

// OutputDateFormat is the format for the date of birth in the response.
const OutputDateFormat = "02.01.2006"

// ValidateByIINResponse is the response structure for the Execute handler.
type ValidateByIINResponse struct {
	Correct     bool   `json:"correct"`
	Sex         string `json:"sex"`
	DateOfBirth string `json:"date_of_birth"`
}

// Execute is a HTTP handler function for validating an IIN.
// It retrieves the IIN from the URL parameter, validates the IIN,
// gets the gender and date of birth from the IIN, and returns a JSON response.
func Execute(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.iin_validate"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		iin := chi.URLParam(r, "iin")
		if iin == "" {
			log.Info("iin is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		err := iin_validator.ValidateIIN(iin)
		if err != nil {
			log.Error("failed to validate IIN", Err(err))
			render.JSON(w, r, ValidateByIINResponse{
				Correct:     false,
				Sex:         "",
				DateOfBirth: "",
			})
			return
		}
		gender, err := iin_validator.GetGender(int(iin[6] - '0'))
		if err != nil {
			log.Error("failed to get gender", Err(err))
			render.JSON(w, r, resp.Error("failed to get gender"))
			return
		}
		dateOfBirth, err := iin_validator.GetDateOfBirth(iin)
		if err != nil {
			log.Error("failed to get gender", Err(err))
			render.JSON(w, r, resp.Error("failed to get gender"))
			return
		}

		formattedDOB := dateOfBirth.Format(OutputDateFormat)
		log.Info("IIN Validated", slog.String("gender", gender), slog.String("date of birth", formattedDOB))
		render.JSON(w, r, ValidateByIINResponse{
			Correct:     true,
			Sex:         gender,
			DateOfBirth: formattedDOB,
		})
	}
}

// Err is a helper function to create a structured log attribute for errors.
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
