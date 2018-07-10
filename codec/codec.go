package codec

type Codec interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) (int, error)
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

func (b *Buffer) Package(i interface{}) error {
	count, err := b.Decode(b.data, i)
	b.Consume(count)
	return err
}
