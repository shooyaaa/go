package session

import (
	"bytes"
	"errors"
	"io"
	"sync"

	"github.com/shooyaaa/core/codec"
	"github.com/shooyaaa/log"
)

type Owner interface {
	OpHandler(codec.Op, *Session)
	SessionClose(int64)
}

var codecInstance codec.Codec[codec.Op]

func SetCodec(c codec.Codec[codec.Op]) {
	codecInstance = c
}

type Session struct {
	Id     int64
	owner  Owner
	Conn   io.ReadWriter
	buffer *bytes.Buffer
}

func (s *Session) WriteWithCodec(msg codec.Op, c codec.Codec[codec.Op]) (int, error) {
	buffer, _ := c.Encode(msg)
	s.Conn.Write(buffer)
	log.DebugF("down write msg")
	return 0, nil
}

func (s *Session) ReadWithCodec(c codec.Codec[codec.Op]) (*codec.Op, error) {
	buffer := make([]byte, 1024)
	count, err := s.Conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	s.buffer.Write(buffer[0:count])
	op, err := c.Decode(s.buffer.Bytes())
	if err != nil {
		return nil, err
	}
	// 清空已读取的数据
	s.buffer.Reset()
	return &op, err
}

func (s *Session) Write(msg codec.Op) (int, error) {
	if codecInstance == nil {
		return -1, errors.New("Default codec should setted")
	}
	return s.WriteWithCodec(msg, codecInstance)
}

func (s *Session) Read() (*codec.Op, error) {
	if codecInstance == nil {
		return nil, errors.New("Default codec should setted")
	}
	return s.ReadWithCodec(codecInstance)
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
			op, err := s.Read()
			if err != nil {
				log.InfoF("error while read from session: %v", err)
				if err == io.EOF {
					log.ErrorF("session end reason %v", err)
					s.owner.SessionClose(s.Id)
					break
				}
			} else {
				s.owner.OpHandler(*op, s)
			}
		}
	})
}
