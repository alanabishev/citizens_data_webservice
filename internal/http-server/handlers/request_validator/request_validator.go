package request_validator

import (
	"citizen_webservice/internal/iin_validator"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func validateIIN(fl validator.FieldLevel) bool {
	iin := fl.Field().String()
	err := iin_validator.ValidateIIN(iin)
	return err == nil
}

func init() {
	validate = validator.New()
	err := validate.RegisterValidation("iin", validateIIN)
	if err != nil {
		panic(err)
	}
}

func GetValidator() *validator.Validate {
	return validate
}
