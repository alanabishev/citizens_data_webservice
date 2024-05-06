package save

import (
	"citizen_webservice/internal/http-server/handlers/request_validator"
	resp "citizen_webservice/internal/http-server/handlers/response"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log/slog"
	"net/http"

	"citizen_webservice/internal/storage"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	IIN   string `json:"iin" validate:"required,len=12,iin"`
	Name  string `json:"name" validate:"required,alphaunicode"`
	Phone string `json:"phone" validate:"required"`
}

type PersonSaver interface {
	SavePerson(iin string, name string, phone string) error
}

func Person(log *slog.Logger, personSaver PersonSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.Person"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		customValidator := request_validator.GetValidator()
		log.Info("request body decoded", slog.Any("request", req))
		if err := customValidator.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", Err(err))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		err = personSaver.SavePerson(req.IIN, req.Name, req.Phone)
		if errors.Is(err, storage.ErrorIINExists) {
			log.Info("iin already exists", slog.String("iin", req.IIN))
			render.JSON(w, r, resp.Error("iin already exists"))
			return
		}
		if errors.Is(err, storage.ErrorPhoneNumberExists) {
			log.Info("phone number already exists", slog.String("phone", req.Phone))
			render.JSON(w, r, resp.Error("phone number already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add person", Err(err))
			render.JSON(w, r, resp.Error("failed to add person"))
			return
		}

		log.Info("person added", slog.String("id", req.IIN))
		render.JSON(w, r, resp.OK())
	}
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
