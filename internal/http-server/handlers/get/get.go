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

type ByIINResponse struct {
	resp.Response
	storage.PersonInfo
}

type ByNameResponse struct {
	resp.Response
	People []storage.PersonInfo `json:"people"`
}

type PersonGetter interface {
	GetPersonByIIN(iin string) (storage.PersonInfo, error)
	GetPersonByName(name string) ([]storage.PersonInfo, error)
}

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
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		err := iin_validator.ValidateIIN(iin)
		if err != nil {
			log.Error("failed to validate IIN", Err(err))
			render.JSON(w, r, resp.Error("failed to validate IIN"))
			return
		}

		personInfo, err := personGetter.GetPersonByIIN(iin)
		if errors.Is(err, storage.ErrorIINNotFound) {
			log.Info("iin not found", slog.String("iin", iin))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, resp.Error("iin not found"))
			return
		}
		if err != nil {
			log.Error("failed to get person", Err(err))
			render.JSON(w, r, resp.Error("failed to get person"))
			return
		}

		log.Info("person retrieved", slog.String("person", fmt.Sprintf("%+v", personInfo)))
		render.JSON(w, r, ByIINResponse{
			Response: resp.OK(),
			PersonInfo: storage.PersonInfo{
				IIN:   personInfo.IIN,
				Name:  personInfo.Name,
				Phone: personInfo.Phone,
			},
		})
	}
}

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
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		peopleInfo, err := personGetter.GetPersonByName(name)
		if errors.Is(err, storage.ErrorNameNotFound) {
			log.Info("name not found", slog.String("name", name))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ByNameResponse{
				Response: resp.OK(),
				People:   []storage.PersonInfo{},
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
			Response: resp.OK(),
			People:   peopleInfo,
		})
	}
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
