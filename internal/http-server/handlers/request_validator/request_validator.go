// Package request_validator provides functionality for validating requests.
package request_validator

import (
	"citizen_webservice/internal/iin_validator"
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// validate is a global instance of the validator.
var validate *validator.Validate

// phoneRegex matches common Kazakhstan phone number formats
var phoneRegex = regexp.MustCompile(`^(\+?7|8)?[0-9]{10}$`)

// validateIIN is a custom validation function for IIN (Individual Identification Number).
// It uses the iin_validator package to validate the IIN.
// It returns true if the IIN is valid, and false otherwise.
func validateIIN(fl validator.FieldLevel) bool {
	iin := fl.Field().String()
	err := iin_validator.ValidateIIN(iin)
	return err == nil
}

// validatePhone is a custom validation function for phone numbers.
// It validates common phone number formats (Kazakhstan format).
// Returns true if the phone number is valid, false otherwise.
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	cleanPhone := regexp.MustCompile(`[\s\-()]`).ReplaceAllString(phone, "")
	return phoneRegex.MatchString(cleanPhone)
}

// validateName is a custom validation function for person names.
// It ensures the name contains only letters, spaces, hyphens, and apostrophes.
// It also checks that the name is between 2 and 100 characters.
// Returns true if the name is valid, false otherwise.
func validateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	if len(name) < 2 || len(name) > 100 {
		return false
	}

	hasLetter := false
	for _, r := range name {
		if unicode.IsLetter(r) {
			hasLetter = true
		} else if r != ' ' && r != '-' && r != '\'' {
			return false
		}
	}

	return hasLetter
}

// init is a special function that is called when the package is initialized.
// It creates a new instance of the validator and registers all custom validation functions.
// If any registration fails, it panics.
func init() {
	validate = validator.New()

	if err := validate.RegisterValidation("iin", validateIIN); err != nil {
		panic(err)
	}

	if err := validate.RegisterValidation("phone", validatePhone); err != nil {
		panic(err)
	}

	if err := validate.RegisterValidation("name", validateName); err != nil {
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
