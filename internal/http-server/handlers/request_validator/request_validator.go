// Package request_validator provides functionality for validating requests.
package request_validator

import (
	"citizen_webservice/internal/iin_validator"
	"github.com/go-playground/validator/v10"
)

// validate is a global instance of the validator.
var validate *validator.Validate

// validateIIN is a custom validation function for IIN (Individual Identification Number).
// It uses the iin_validator package to validate the IIN.
// It returns true if the IIN is valid, and false otherwise.
func validateIIN(fl validator.FieldLevel) bool {
	iin := fl.Field().String()
	err := iin_validator.ValidateIIN(iin)
	return err == nil
}

// init is a special function that is called when the package is initialized.
// It creates a new instance of the validator and registers the custom IIN validation function.
// If the registration fails, it panics.
func init() {
	validate = validator.New()
	err := validate.RegisterValidation("iin", validateIIN)
	if err != nil {
		panic(err)
	}
}

// GetValidator is a function that returns the global instance of the validator.
func GetValidator() *validator.Validate {
	return validate
}

// CheckErrorIsValidation is a function that checks if an error is a validation error.
// It returns true if the error is a validation error, and false otherwise.
func CheckErrorIsValidation(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}
