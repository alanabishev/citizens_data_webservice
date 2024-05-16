package iin_validator

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetGender(t *testing.T) {
	testCases := []struct {
		name     string
		digit    int
		expected string
		err      error
	}{
		{
			name:     "Test Case 1: 19th Century Male",
			digit:    SeventhDigitMaleIIN19Century,
			expected: "male",
			err:      nil,
		},
		{
			name:     "Test Case 2: 20th Century Female",
			digit:    SeventhDigitFemaleIIN20Century,
			expected: "female",
			err:      nil,
		},
		{
			name:     "Test Case 3: Invalid Digit",
			digit:    0,
			expected: "",
			err:      fmt.Errorf("invalid digit, must be between %d and %d inclusive", SeventhDigitMaleIIN19Century, SeventhDigitFemaleIIN21Century),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GetGender(tc.digit)
			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestGetDateOfBirth(t *testing.T) {
	testCases := []struct {
		name     string
		iin      string
		expected time.Time
		err      error
	}{
		{
			name:     "Test Case 1: Valid IIN",
			iin:      "830218350074",
			expected: time.Date(1983, 2, 18, 0, 0, 0, 0, time.UTC),
			err:      nil,
		},
		{
			name:     "Test Case 2: IIN too short",
			iin:      "12345",
			expected: time.Time{},
			err:      fmt.Errorf("input string is too short"),
		},
		{
			name:     "Test Case 3: Invalid Date",
			iin:      "990230350074",
			expected: time.Time{},
			err:      fmt.Errorf("parsed date format is incorrect: %w", &time.ParseError{Layout: "20060102", Value: "19990230", Message: ": day out of range"}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GetDateOfBirth(tc.iin)
			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestGetCenturyOfBirth(t *testing.T) {
	testCases := []struct {
		name     string
		digit    int
		expected int
		err      error
	}{
		{
			name:     "Test Case 1: 19th Century Male",
			digit:    SeventhDigitMaleIIN19Century,
			expected: 19,
			err:      nil,
		},
		{
			name:     "Test Case 2: 20th Century Female",
			digit:    SeventhDigitFemaleIIN20Century,
			expected: 20,
			err:      nil,
		},
		{
			name:     "Test Case 3: Invalid Digit",
			digit:    -1,
			expected: 0,
			err:      fmt.Errorf("invalid digit, must be between %d and %d inclusive", SeventhDigitMaleIIN19Century, SeventhDigitFemaleIIN21Century),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := getCenturyOfBirth(tc.digit)
			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestValidateSeventhDigit(t *testing.T) {
	testCases := []struct {
		name  string
		digit int
		err   error
	}{
		{
			name:  "Test Case 1: Valid Digit",
			digit: SeventhDigitMaleIIN19Century,
			err:   nil,
		},
		{
			name:  "Test Case 2: Invalid Digit",
			digit: 0,
			err:   fmt.Errorf("invalid digit, must be between %d and %d inclusive", SeventhDigitMaleIIN19Century, SeventhDigitFemaleIIN21Century),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateSeventhDigit(tc.digit)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestCalculate12thDigit(t *testing.T) {
	testCases := []struct {
		name      string
		iin       string
		algorithm int
		expected  int
		err       error
	}{
		{
			name:      "Test Case 1: Valid IIN with Algorithm 1",
			iin:       "123456789012",
			algorithm: Algorithm1,
			expected:  10,
			err:       nil,
		},
		{
			name:      "Test Case 2: Valid IIN with Algorithm 2",
			iin:       "123456789012",
			algorithm: Algorithm2,
			expected:  3,
			err:       nil,
		},
		{
			name:      "Test Case 3: IIN too short",
			iin:       "12345",
			algorithm: Algorithm1,
			expected:  0,
			err:       fmt.Errorf("iin string too short"),
		},
		{
			name:      "Test Case 4: Valid IIN",
			iin:       "830218350074",
			algorithm: Algorithm1,
			expected:  4,
			err:       nil,
		},
		{
			name:      "Test Case 5: Valid IIN, Algorithm 1",
			iin:       "600426400918",
			algorithm: Algorithm1,
			expected:  10,
			err:       nil,
		},
		{
			name:      "Test Case 5: Valid IIN Algorithm 2",
			iin:       "600426400918",
			algorithm: Algorithm2,
			expected:  8,
			err:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := calculate12thDigit(tc.iin, tc.algorithm)
			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.err, err)
		})
	}
}
