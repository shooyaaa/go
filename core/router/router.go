package router

import (
	"errors"
	"fmt"
	"github.com/shooyaaa/config"
	"github.com/shooyaaa/core"
	"github.com/shooyaaa/core/library"
	"github.com/shooyaaa/core/network"
	"github.com/shooyaaa/log"
	"strconv"
	"strings"
)

type ROUTER_TYPE int

const (
	TCP_ROUTER = iota
	CHAN_ROUTER
)

var CENTRAL_REGISTRY string = TcpRouterName(config.CentralRegistryName)

func TcpRouterName(name string) string {
	return fmt.Sprintf("%v:%v", TCP_ROUTER, name)
}

type Router interface {
	Forward(*Package) error
	ToString() string
}

type TcpRouter struct {
	host string
	port int
}

func (tr *TcpRouter) Forward(p *Package) error {
	log.DebugF("tcp router forward package seq: %v, ack: %v", p.seq, p.ack)
	conn := network.TcpConn{}
	err := conn.Dial(tr.host, tr.port)
	if err != nil {
		return err
	}
	conn.Write(p.body)
	return nil
}

func (tr *TcpRouter) ToString() string {
	return fmt.Sprintf("%v:%v:%v", TCP_ROUTER, tr.host, tr.port)
}
func NewTcpRouter(host string, port int) *TcpRouter {
	tr := TcpRouter{}
	tr.host = host
	tr.port = port
	return &tr
}

type ChanRouter struct {
	ch chan *Package
}

func (cr *ChanRouter) Forward(p *Package) error {
	cr.ch <- p
	return nil
}

func (cr *ChanRouter) ToString() string {
	return fmt.Sprintf("%v:%v", CHAN_ROUTER, cr.ch)
}
func NewChanRouter(ch chan *Package) *ChanRouter {
	cr := ChanRouter{}
	cr.ch = ch
	return &cr
}

var dummyRegistry DummyRegistry

var tcpRegistry TcpRegistry

func LookUp(entity string) (Router, error) {
	info := strings.Split(entity, ":")
	routerType, err := strconv.Atoi(info[0])
	if err != nil {
		return nil, err
	}
	if len(info) < 2 {
		return nil, errors.New(core.PARAMS_ERROR)
	}
	switch routerType {
	case TCP_ROUTER:
		addr, err := tcpRegistry.Get(info[1])
		if err != nil {
			return nil, err
		}
		data := strings.Split(addr, ":")
		port, err := strconv.Atoi(data[1])
		if err != nil {
			return nil, err
		}
		return NewTcpRouter(data[0], port), nil
	case CHAN_ROUTER:
		ch, err := dummyRegistry.Get(info[1])
		if err != nil {
			return nil, err
		}
		return NewChanRouter(ch), nil
	}
	return nil, errors.New(core.UNKNOWN)
}

func init() {
	dummyRegistry = DummyRegistry{tables: map[string]chan *Package{}}
	tcpRegistry = TcpRegistry{tables: library.Redis{}}
	tcpRegistry.tables.Init(map[string]interface{}{"address": config.RegistryRedisAddress})
}

func AddTcpAddress(name string, addr string) {
	tcpRegistry.Set(name, addr)
}

func AddDummyAddress(name string, ch chan *Package) {
	dummyRegistry.Set(name, ch)
}
