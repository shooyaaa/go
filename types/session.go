package types

import (
	"errors"
	"io"
)

type Handler interface {
}

var codec Codec

func SetCodec(c Codec) {
	codec = c
}

type Session struct {
	Id      int64
	handler Handler
	Conn    io.ReadWriter
}

func (s *Session) WriteWithCodec(msg interface{}, c Codec) (int, error) {

	return 0, nil
}

func (s *Session) ReadWithCodec(c Codec) ([]byte, error) {

	return nil, nil
}

func (s *Session) Write(msg interface{}) (int, error) {
	if codec == nil {
		return -1, errors.New("Default codec should setted")
	}
	return s.WriteWithCodec(msg, codec)
}

func (s *Session) Read() ([]byte, error) {
	if codec == nil {
		return nil, errors.New("Default codec should setted")
	}
	return s.ReadWithCodec(codec)
}
