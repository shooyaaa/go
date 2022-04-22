package core

import (
	"errors"
	"fmt"
)

type error_code int

const (
	NOT_FOUND    = "NOT FOUND"
	PARAMS_ERROR = "PARAMS ERROR"
	TYPE_ERROR   = "TYPE_ERROR"
	UNKNOWN      = "UNKNOWN_ERROR"
)

func New(msg string, code error_code) error {
	return errors.New(fmt.Sprintf("code: %v, msg: %v", code, msg))
}
