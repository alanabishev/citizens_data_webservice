package iin_validator

import (
	"fmt"
	"time"
	"unicode"
)

const (
	IINLength                      = 12
	FirstSixDigitsDateFormat       = "20060102" // YYYYMMDD
	SeventhDigitMaleIIN19Century   = 1
	SeventhDigitMaleIIN20Century   = 3
	SeventhDigitMaleIIN21Century   = 5
	SeventhDigitFemaleIIN19Century = 2
	SeventhDigitFemaleIIN20Century = 4
	SeventhDigitFemaleIIN21Century = 6
	TwelfthDigitForSecondAlgorithm = 10
	Algorithm1                     = 1
	Algorithm2                     = 2
)

var (
	weightsAlgorithm1 = [11]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	weightsAlgorithm2 = [11]int{3, 4, 5, 6, 7, 8, 9, 10, 11, 1, 2}
)

func ValidateIIN(iin string) error {
	if len(iin) != IINLength {
		return fmt.Errorf("IIN must be %d digits long", IINLength)
	}
	for _, r := range iin {
		if !unicode.IsDigit(r) {
			return fmt.Errorf("IIN must only contain digits")
		}
	}
	_, err := GetDateOfBirth(iin)
	if err != nil {
		return fmt.Errorf("first 6 digits of IIN must form a valid date (YYMMDD): %w", err)
	}
	_, err = GetGender(int(iin[6] - '0'))
	if err != nil {
		return fmt.Errorf("invalid 7th digit: %w", err)
	}

	twelfthDigit, err := calculate12thDigit(iin, Algorithm1)
	if err != nil {
		return fmt.Errorf("error while validating 12th digit: %w", err)
	}
	if twelfthDigit == TwelfthDigitForSecondAlgorithm {
		twelfthDigit, err = calculate12thDigit(iin, Algorithm2)
		if err != nil {
			return fmt.Errorf("error while validating 12th digit: %w", err)
		}
	}
	if twelfthDigit != int(iin[11]-'0') {
		return fmt.Errorf("invalid 12th digit")
	}

	return nil
}

func GetGender(digit int) (string, error) {
	if err := validateSeventhDigit(digit); err != nil {
		return "", err
	}
	switch digit {
	case SeventhDigitFemaleIIN19Century, SeventhDigitFemaleIIN20Century, SeventhDigitFemaleIIN21Century:
		return "female", nil
	case SeventhDigitMaleIIN19Century, SeventhDigitMaleIIN20Century, SeventhDigitMaleIIN21Century:
		return "male", nil
	}
	return "", fmt.Errorf("invalid digit")
}

func GetDateOfBirth(iin string) (time.Time, error) {
	if len(iin) < 7 {
		return time.Time{}, fmt.Errorf("input string is too short")
	}
	centuryOfBirth, err := getCenturyOfBirth(int(iin[6] - '0'))
	var date time.Time

	if err != nil {
		return date, fmt.Errorf("error while getting century of birth: %w", err)
	}
	dobToParse := fmt.Sprintf("%d%s", centuryOfBirth-1, iin[:6])
	date, err = time.Parse(FirstSixDigitsDateFormat, dobToParse)
	if err != nil || date.After(time.Now()) {
		return date, fmt.Errorf("parsed date format is incorrect: %w", err)
	}
	return date, nil
}

func getCenturyOfBirth(digit int) (int, error) {
	if err := validateSeventhDigit(digit); err != nil {
		return 0, err
	}
	switch digit {
	case SeventhDigitMaleIIN19Century, SeventhDigitFemaleIIN19Century:
		return 19, nil
	case SeventhDigitMaleIIN20Century, SeventhDigitFemaleIIN20Century:
		return 20, nil
	case SeventhDigitMaleIIN21Century, SeventhDigitFemaleIIN21Century:
		return 21, nil
	}
	return 0, fmt.Errorf("invalid digit")
}

func validateSeventhDigit(digit int) error {
	if digit < SeventhDigitMaleIIN19Century || digit > SeventhDigitFemaleIIN21Century {
		return fmt.Errorf("invalid digit, must be between %d and %d inclusive", SeventhDigitMaleIIN19Century, SeventhDigitFemaleIIN21Century)
	}
	return nil
}

func calculate12thDigit(iin string, algorithm int) (int, error) {
	if len(iin) < 11 {
		return 0, fmt.Errorf("iin string too short")
	}
	weights := weightsAlgorithm1
	if algorithm == Algorithm2 {
		weights = weightsAlgorithm2
	}

	sum := 0
	for i := 0; i < 11; i++ {
		digit := int(iin[i] - '0')
		sum += weights[i] * digit
	}

	digit12 := sum % 11
	return digit12, nil
}
