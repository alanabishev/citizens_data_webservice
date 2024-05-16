package tests

import (
	"fmt"
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"net/url"
	"testing"
)

const (
	host = "0.0.0.0:8082"
)

func TestResponse_ValidIIN(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())
	iinList := []string{"600426400918", "790708301327", "990109300285"}
	sexList := []string{"female", "male", "male"}
	dobList := []string{"26.04.1960", "08.07.1979", "09.01.1999"}
	for idx := 0; idx < len(iinList); idx++ {
		e.GET(fmt.Sprintf("/iin_check/%s", iinList[idx])).
			WithBasicAuth("user", "password").
			Expect().
			Status(http.StatusOK).
			JSON().Object().
			ContainsKey("correct").HasValue("correct", true).
			ContainsKey("sex").HasValue("sex", sexList[idx]).
			ContainsKey("date_of_birth").HasValue("date_of_birth", dobList[idx])
	}
}

func TestResponse_InvalidIIN(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	iinList := []string{"600426400919", "600", "asdasd"}
	for _, iin := range iinList {
		e.GET(fmt.Sprintf("/iin_check/%s", iin)).
			WithBasicAuth("user", "password").
			Expect().
			Status(http.StatusOK).
			JSON().Object().
			ContainsKey("correct").HasValue("correct", false).
			ContainsKey("sex").HasValue("sex", "").
			ContainsKey("date_of_birth").HasValue("date_of_birth", "")
	}

}

func TestSavePersonEndpoint(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	test_iin := "980301450725"

	// 1) Empty JSON body
	e.POST("/people/info").
		WithBasicAuth("user", "password").
		WithJSON(map[string]interface{}{}).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ContainsKey("success").HasValue("success", false).
		ContainsKey("errors").NotEmpty()

	// 2) Invalid IIN
	e.POST("/people/info").
		WithBasicAuth("user", "password").
		WithJSON(map[string]interface{}{
			"iin":   "1234",
			"name":  "Test Name",
			"phone": "1234567890",
		}).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ContainsKey("success").HasValue("success", false).
		ContainsKey("errors").NotEmpty()

	// 3) Valid IIN
	e.POST("/people/info").
		WithBasicAuth("user", "password").
		WithJSON(map[string]interface{}{
			"iin":   test_iin,
			"name":  "Test Name",
			"phone": "1234567890",
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("success").HasValue("success", true).
		ContainsKey("errors").HasValue("errors", nil)

	// 4) The same person, but this time expect an error of iin exists
	e.POST("/people/info").
		WithBasicAuth("user", "password").
		WithJSON(map[string]interface{}{
			"iin":   test_iin,
			"name":  "Test Name",
			"phone": "1234567891",
		}).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().
		ContainsKey("success").HasValue("success", false).
		ContainsKey("errors").HasValue("errors", []string{"Failed to save person: storage.sqlite.SavePerson: IIN already exists"})

	//5) Different valid IIN, but the phone number is the same
	e.POST("/people/info").
		WithBasicAuth("user", "password").
		WithJSON(map[string]interface{}{
			"iin":   "600426400918",
			"name":  "Test Name",
			"phone": "1234567890",
		}).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().
		ContainsKey("success").HasValue("success", false).
		ContainsKey("errors").HasValue("errors", []string{"Failed to save person: storage.sqlite.SavePerson: phone number already exists"})

	// Delete a person with a specific IIN
	e.DELETE(fmt.Sprintf("/people/delete/%s", test_iin)).
		WithBasicAuth("user", "password").
		Expect().
		Status(http.StatusOK)
}

func TestGetPersonByIINEndpoint(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	test_iin := "980301450725"

	// 1) Get a person with empty IIN parameter
	e.GET("/people/info/iin/").
		WithBasicAuth("user", "password").
		Expect().
		Status(http.StatusNotFound)

	// 2) Get a person with invalid IIN
	e.GET("/people/info/iin/1234").
		WithBasicAuth("user", "password").
		Expect().
		Status(http.StatusBadRequest)

	// 3) Get a person with valid IIN but the person is not yet in the database
	e.GET(fmt.Sprintf("/people/info/iin/%s", test_iin)).
		WithBasicAuth("user", "password").
		Expect().
		Status(http.StatusNotFound)

	// 4) Create a person using the POST /people/info endpoint
	e.POST("/people/info").
		WithBasicAuth("user", "password").
		WithJSON(map[string]interface{}{
			"iin":   test_iin,
			"name":  "Test Name",
			"phone": "1234567890",
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("success").HasValue("success", true).
		ContainsKey("errors").HasValue("errors", nil)

	e.GET(fmt.Sprintf("/people/info/iin/%s", test_iin)).
		WithBasicAuth("user", "password").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("success").HasValue("success", true).
		ContainsKey("IIN").HasValue("IIN", test_iin)

	// And delete him
	e.DELETE(fmt.Sprintf("/people/delete/%s", test_iin)).
		WithBasicAuth("user", "password").
		Expect().
		Status(http.StatusOK)
}

func TestGetPersonByNameEndpoint(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	test_name := "Sally"

	// 1) Get a person with empty name parameter
	e.GET("/people/info/name/").
		WithBasicAuth("user", "password").
		Expect().
		Status(http.StatusNotFound)

	// 2) Create a person using the POST /people/info endpoint, and get him
	e.POST("/people/info").
		WithBasicAuth("user", "password").
		WithJSON(map[string]interface{}{
			"iin":   "980301450725",
			"name":  test_name,
			"phone": "1234567890",
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("success").HasValue("success", true).
		ContainsKey("errors").HasValue("errors", nil)

	e.GET(fmt.Sprintf("/people/info/name/%s", "ll")).
		WithBasicAuth("user", "password").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("success").HasValue("success", true).
		ContainsKey("people").NotEmpty()

	// 3) Create another person with similar name and get by name
	e.POST("/people/info").
		WithBasicAuth("user", "password").
		WithJSON(map[string]interface{}{
			"iin":   "790708301327",
			"name":  "Lilly",
			"phone": "1234567891",
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("success").HasValue("success", true).
		ContainsKey("errors").HasValue("errors", nil)

	e.GET(fmt.Sprintf("/people/info/name/%s", "l")).
		WithBasicAuth("user", "password").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("success").HasValue("success", true).
		ContainsKey("people").NotEmpty()

	// 5) Get name with symbols not used in previous, assert the result array is empty
	e.GET("/people/info/name/qqqq").
		WithBasicAuth("user", "password").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("success").HasValue("success", true).
		ContainsKey("people").HasValue("people", nil)

	// 6) Delete the 2 people
	e.DELETE("/people/delete/790708301327").
		WithBasicAuth("user", "password").
		Expect().
		Status(http.StatusOK)

	e.DELETE("/people/delete/980301450725").
		WithBasicAuth("user", "password").
		Expect().
		Status(http.StatusOK)
}
