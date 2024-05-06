package tests

import (
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"net/url"
	"testing"
)

const (
	host = "localhost:8082"
)

func TestResponse_ValidIIN(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	e.GET("/iin_check/600426400918").
		WithBasicAuth("user", "password").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("status").HasValue("status", "OK").
		ContainsKey("correct").HasValue("correct", true).
		ContainsKey("sex").HasValue("sex", "female").
		ContainsKey("date_of_birth").HasValue("date_of_birth", "26.04.1960")
}

func TestResponse_InvalidIIN(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	e.GET("/iin_check/600426400919").
		WithBasicAuth("user", "password").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("status").HasValue("status", "OK").
		ContainsKey("correct").HasValue("correct", false).
		ContainsKey("sex").HasValue("sex", "").
		ContainsKey("date_of_birth").HasValue("date_of_birth", "")
}
