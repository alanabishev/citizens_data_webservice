package storage

import "errors"

var (
	ErrorIINNotFound       = errors.New("IIN not found")
	ErrorIINExists         = errors.New("IIN already exists")
	ErrorNameNotFound      = errors.New("name not found")
	ErrorPhoneNumberExists = errors.New("phone number already exists")
)

type PersonInfo struct {
	IIN   string
	Name  string
	Phone string
}
