package types

import (
	"log"
)

type Codec interface {
	Encode([]Op) ([]byte, error)
	Decode([]byte, *[]Op) (int, error)
}

type Buffer struct {
	data []byte
	Codec
}

func (b *Buffer) Append(bt []byte) (int, error) {
	b.data = append(b.data, bt...)
	return len(bt), nil
}

func (b *Buffer) Consume(i int) (int, error) {
	b.data = b.data[i:]
	return i, nil
}

func (b *Buffer) Package(data []byte) ([]Op, error) {
	//op := binary.BigEndian.Uint16(data)
	//dispatcher.Handle(OpCode(op), data[2:])
	//b.Consume(len(b.data))
	ops := make([]Op, 1)
	size, err := b.Decode(data, &ops)
	if err != nil {
		log.Printf("Error while decode buffer %v", err)
	}
	data = data[size:]
	return ops, err
}
