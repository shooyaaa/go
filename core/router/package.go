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
	codec := codec.GetCodec(p.codec)
	if codec == nil {
		panic("package has invalid code type")
	}
	body, err := codec.Encode(p.body)
	if err != nil {
		return err
	}
	p.body = body
	return nil
}

func (p *Package) GetBody() []byte {
	return p.body
}

func (p *Package) Decode() (interface{}, error) {
	codec := codec.GetCodec(p.codec)
	ret, _, err := codec.Decode(p.body)
	if err != nil {
		return nil, err
	}
	return ret, err
}
