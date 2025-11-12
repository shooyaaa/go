package router

import (
	"github.com/shooyaaa/core/codec"
)

type PackageHandler interface {
	HandlePackage(*Package) error
}

type Header struct {
	seq   int64
	ack   int64
	len   int
	codec codec.CODEC_TYPE
}

type Package struct {
	Header
	body []byte
}

func (p *Package) Encode() error {
	// Encode is a no-op since body is already []byte
	// If encoding is needed, it should be done before setting the body
	return nil
}

func (p *Package) GetBody() []byte {
	return p.body
}

func (p *Package) Decode() (interface{}, error) {
	codecInstance := codec.NewCodec[codec.Op](p.codec)
	ret, err := codecInstance.Decode(p.body)
	if err != nil {
		return nil, err
	}
	return ret, err
}
