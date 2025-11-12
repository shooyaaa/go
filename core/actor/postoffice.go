package actor

import (
	"context"

	"github.com/shooyaaa/core"
	"github.com/shooyaaa/core/library"
)

type Postoffice interface {
	Add(ctx context.Context, a PostManAddress) *core.CoreError
	Remove(ctx context.Context, a PostManAddress) *core.CoreError
	Dispatch(ctx context.Context, mail Mail[any]) *core.CoreError
}

type postofficeImpl struct {
	h library.ConsistentHash[PostManAddress]
}

func NewPostoffice(h library.ConsistentHash[PostManAddress]) Postoffice {
	return &postofficeImpl{h: h}
}

func (p *postofficeImpl) Add(ctx context.Context, a PostManAddress) *core.CoreError {
	p.h.Add(a)
	return nil
}

func (p *postofficeImpl) Remove(ctx context.Context, a PostManAddress) *core.CoreError {
	p.h.Remove(a)
	return nil
}

func (p *postofficeImpl) Dispatch(ctx context.Context, mail Mail[any]) *core.CoreError {
	receiver := mail.Receiver()
	a, ok := p.h.Get((&receiver).String())
	if ok {
		return a.Transform(ctx, mail)
	}
	return core.NewCoreError(core.ERROR_CODE_POSTMAN_NOT_FOUND, "postman not found in dispatch")
}
