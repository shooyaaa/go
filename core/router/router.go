package router

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shooyaaa/config"
	"github.com/shooyaaa/core/network"
	"github.com/shooyaaa/core/storage"
	"github.com/shooyaaa/log"
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
	tcp  network.Tcp
	TcpRegistry
}

func (tr *TcpRouter) Forward(p *Package) error {
	log.DebugF("tcp router forward package seq: %v, ack: %v", p.seq, p.ack)
	conn := network.TcpConn{}
	_, err := conn.Dial(tr.host, tr.port)
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

func (tr *TcpRouter) Listen(addr string) error {
	return tr.tcp.Listen(addr)
}
func NewTcpRouter(host string, port int, cache storage.Cache) *TcpRouter {
	tr := TcpRouter{}
	tr.host = host
	tr.port = port
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

var globalTcpRegistry *TcpRegistry

func init() {
	//tcpRegistry.tables.Init(map[string]interface{}{"address": config.RegistryRedisAddress})
	globalTcpRegistry = &TcpRegistry{}
}

func LookUp(routerId string) (Router, error) {
	parts := strings.Split(routerId, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid router id format: %s", routerId)
	}
	routerType, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	entity := parts[1]

	switch ROUTER_TYPE(routerType) {
	case TCP_ROUTER:
		addr, err := globalTcpRegistry.Get(entity)
		if err != nil {
			return nil, err
		}
		data := strings.Split(addr, ":")
		if len(data) != 2 {
			return nil, fmt.Errorf("invalid address format: %s", addr)
		}
		port, err := strconv.Atoi(data[1])
		if err != nil {
			return nil, err
		}
		return NewTcpRouter(data[0], port, globalTcpRegistry.tables), nil
	default:
		return nil, fmt.Errorf("unknown router type: %d", routerType)
	}
}

func AddTcpAddress(name string, address string) error {
	return globalTcpRegistry.Set(name, address)
}
