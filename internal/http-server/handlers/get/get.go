// Package get provides HTTP handlers for getting person information.
package get

import (
	resp "citizen_webservice/internal/http-server/handlers/response"
	"citizen_webservice/internal/iin_validator"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"

	"citizen_webservice/internal/storage"
	"github.com/go-chi/render"
)

// ByIINResponse is the response structure for the ByIIN handler.
type ByIINResponse struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
	storage.PersonInfo
}

// ByNameResponse is the response structure for the ByName handler.
type ByNameResponse struct {
	Success bool                 `json:"success"`
	Errors  []string             `json:"errors"`
	People  []storage.PersonInfo `json:"people"`
}

// PersonGetter is an interface for getting person information.
type PersonGetter interface {
	GetPersonByIIN(iin string) (storage.PersonInfo, error)
	GetPersonByName(name string) ([]storage.PersonInfo, error)
}

// ByIIN is a HTTP handler function for getting a person by their IIN.
// It validates the IIN, retrieves the person information from the storage,
// and returns a JSON response.
func ByIIN(log *slog.Logger, personGetter PersonGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get.ByIIN"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		iin := chi.URLParam(r, "iin")
		if iin == "" {
			log.Info("iin is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ByIINResponse{
				Success: false,
				Errors:  []string{"invalid request, iin is empty"},
			})
			return
		}
		err := iin_validator.ValidateIIN(iin)
		if err != nil {
			log.Error("failed to validate IIN", Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ByIINResponse{
				Success: false,
				Errors:  []string{"failed to validate IIN"},
			})
			return
		}

		personInfo, err := personGetter.GetPersonByIIN(iin)
		if errors.Is(err, storage.ErrorIINNotFound) {
			log.Info("iin not found", slog.String("iin", iin))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ByIINResponse{
				Success: false,
				Errors:  []string{"iin not found"},
			})
			return
		}
		if err != nil {
			log.Error("failed to get person", Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ByIINResponse{
				Success: false,
				Errors:  []string{"failed to get person"},
			})
			return
		}

		log.Info("person retrieved", slog.String("person", fmt.Sprintf("%+v", personInfo)))
		render.JSON(w, r, ByIINResponse{
			Success: true,
			PersonInfo: storage.PersonInfo{
				IIN:   personInfo.IIN,
				Name:  personInfo.Name,
				Phone: personInfo.Phone,
			},
		})
	}
}

// ByName is a HTTP handler function for getting persons by their name.
// It retrieves the person information from the storage,
// and returns a JSON response.
func ByName(log *slog.Logger, personGetter PersonGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get.ByName"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		name := chi.URLParam(r, "name")
		if name == "" {
			log.Info("name is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		peopleInfo, err := personGetter.GetPersonByName(name)
		if errors.Is(err, storage.ErrorNameNotFound) {
			log.Info("name not found", slog.String("name", name))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ByNameResponse{
				People: []storage.PersonInfo{},
			})
			return
		}
		if err != nil {
			log.Error("failed to retrieve people", Err(err))
			render.JSON(w, r, resp.Error("failed to retrieve people"))
			return
		}

		log.Info("person match success", slog.String("matches", fmt.Sprintf("%+v", peopleInfo)))
		render.JSON(w, r, ByNameResponse{
			Success: true,
			People:  peopleInfo,
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
