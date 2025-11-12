package core

import (
	"fmt"
)

type errorCode int

const (
	ERROR_CODE_ACTOR_NOT_FOUND errorCode = iota + 1
	ERROR_CODE_CODEC_ENCODE_ERROR
	ERROR_CODE_CODEC_DECODE_ERROR
	ERROR_CODE_ACTOR_NOT_FOUND_IN_POSTMAN
	ERROR_CODE_POSTMAN_NOT_FOUND
	ERROR_CODE_POSTOFFICE_NOT_REGISTERED
	ERROR_CODE_ACTOR_ALREADY_EXISTS
	ERROR_CODE_MAILBOX_SEND_ERROR
	ERROR_CODE_MAILBOX_RECEIVE_ERROR
	ERROR_TYPE_CORE
	ERROR_TYPE_BUSINESS
)

type CoreError struct {
	code errorCode
	msg  string
}

func NewCoreError(code errorCode, msg string) *CoreError {
	return &CoreError{code: code, msg: msg}
}

func (e *CoreError) String() string {
	return fmt.Sprintf("code: %d, msg: %s", e.code, e.msg)
}

func (e *CoreError) Code() errorCode {
	return e.code
}

func (e *CoreError) Msg() string {
	return e.msg
}
