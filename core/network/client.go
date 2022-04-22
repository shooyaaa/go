package network

import (
	"github.com/shooyaaa/core/session"
	"net/rpc"
)

type Client interface {
	Dial(addr string) session.Session
}

type RpcClient interface {
	Call(addr string, request rpc.Request)
}
