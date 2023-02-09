package app

import (
	"errors"
	"reflect"
)

type ServerError struct {
	error
}

func UCaseError(code string) ServerError {
	return ServerError{
		error: errors.New(code),
	}
}

func IsServerError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(ServerError{})
}
