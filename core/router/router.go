package router

import (
	"fmt"
	"github.com/shooyaaa/config"
	"github.com/shooyaaa/core/network"
	"github.com/shooyaaa/core/storage"
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
	LookUp(id string) (Router, error)
}

type TcpRouter struct {
	host string
	port int
	TcpRegistry
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

func (tr *TcpRouter) LookUp(entity string) (Router, error) {
	addr, err := tr.Get(entity)
	if err != nil {
		return nil, err
	}
	data := strings.Split(addr, ":")
	port, err := strconv.Atoi(data[1])
	if err != nil {
		return nil, err
	}
	return NewTcpRouter(data[0], port, tr.tables), nil
}

func (tr *TcpRouter) ToString() string {
	return fmt.Sprintf("%v:%v:%v", TCP_ROUTER, tr.host, tr.port)
}
func NewTcpRouter(host string, port int, cache storage.Cache) *TcpRouter {
	tr := TcpRouter{}
	tr.tables = cache
	return &tr
}

type ChanRouter struct {
	ch chan *Package
	DummyRegistry
}

func (cr *ChanRouter) Forward(p *Package) error {
	cr.ch <- p
	return nil
}

func (this *ChanRouter) LookUp(entity string) (Router, error) {
	ch, err := this.Get(entity)
	if err != nil {
		return nil, err
	}
	return NewChanRouter(ch), nil
}

func (cr *ChanRouter) ToString() string {
	return fmt.Sprintf("%v:%v", CHAN_ROUTER, cr.ch)
}
func NewChanRouter(ch chan *Package) *ChanRouter {
	cr := ChanRouter{}
	cr.ch = ch
	return &cr
}

func init() {
	//tcpRegistry.tables.Init(map[string]interface{}{"address": config.RegistryRedisAddress})
}
