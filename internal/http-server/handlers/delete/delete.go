// Package delete provides HTTP handlers for deleting person information.
package delete

import (
	resp "citizen_webservice/internal/http-server/handlers/response"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"

	"citizen_webservice/internal/storage"
	"github.com/go-chi/render"
)

// PersonDeleter is an interface for deleting person information.
type PersonDeleter interface {
	DeletePersonByIIN(iin string) error
}

// ByIIN is an HTTP handler function for deleting a person by their IIN.
// It retrieves the IIN from the URL parameter, deletes the person from the storage,
// and returns a JSON response.
func ByIIN(log *slog.Logger, personDeleter PersonDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.ByIIN"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		iin := chi.URLParam(r, "iin")
		if iin == "" {
			log.Info("iin is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("iin is empty"))
			return
		}

		err := personDeleter.DeletePersonByIIN(iin)
		if errors.Is(err, storage.ErrorIINNotFound) {
			log.Info("iin not found", slog.String("iin", iin))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, resp.Error("iin not found"))
			return
		}
		if err != nil {
			log.Error("failed to delete person", Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete person"))
			return
		}

		log.Info("person deleted", slog.String("iin", iin))
		render.JSON(w, r, resp.OK())
	}
}

// Err is a helper function to create a structured log attribute for errors.
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
