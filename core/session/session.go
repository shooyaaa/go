package session

import (
	"bytes"
	"errors"
	codec2 "github.com/shooyaaa/core/codec"
	"github.com/shooyaaa/core/op"
	"github.com/shooyaaa/log"
	"io"
	"sync"
)

type Owner interface {
	OpHandler(op.Op, *Session)
	SessionClose(int64)
}

var codec codec2.Codec

func SetCodec(c codec2.Codec) {
	codec = c
}

type Session struct {
	Id     int64
	owner  Owner
	Conn   io.ReadWriter
	buffer *bytes.Buffer
}

func (s *Session) WriteWithCodec(msg []op.Op, c codec2.Codec) (int, error) {
	buffer, _ := c.Encode(msg)
	s.Conn.Write(buffer)
	return 0, nil
}

func (s *Session) ReadWithCodec(c codec2.Codec) ([]op.Op, error) {
	buffer := make([]byte, 1024)
	count, err := s.Conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	s.buffer.Write(buffer[0:count])
	ops := make([]op.Op, 0)
	i, reduced, err := c.Decode(s.buffer.Bytes())
	ops = i.([]op.Op)
	s.buffer.Next(reduced)
	return ops, err
}

func (s *Session) Write(msg []op.Op) (int, error) {
	if codec == nil {
		return -1, errors.New("Default codec should setted")
	}
	return s.WriteWithCodec(msg, codec)
}

func (s *Session) Read() ([]op.Op, error) {
	if codec == nil {
		return nil, errors.New("Default codec should setted")
	}
	return s.ReadWithCodec(codec)
}

func (s *Session) SetOwner(o Owner) {
	s.owner = o
	var once sync.Once
	go once.Do(func() {
		for {
			if s.buffer == nil {
				s.buffer = &bytes.Buffer{}
			}
			log.DebugF("read from session %d", s.Id)
			ops, err := s.Read()
			if err != nil {
				log.ErrorF("error while read from session: %v", err)
				if err == io.EOF {
					log.ErrorF("session end reason %v", err)
					s.owner.SessionClose(s.Id)
					break
				}
			} else {
				for _, op := range ops {
					s.owner.OpHandler(op, s)
				}
			}
		}
	})
}
