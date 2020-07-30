package connector

import (
	"log"
	"net"
	"time"

	"github.com/shooyaaa/manager"
	"github.com/shooyaaa/types"
)

type Tcp struct {
	Id        types.UUID
	Sessions  map[types.ID]types.Session
	HeartBeat time.Duration
}

type TcpConn struct {
	conn net.Conn
}

func (tc TcpConn) Read() ([]byte, error) {
	bytes := make([]byte, 1024)
	_, err := tc.conn.Read(bytes)
	return bytes, err
}

func (tc TcpConn) Write(bytes []byte) (int, error) {
	return tc.conn.Write(bytes)
}

func (tcp *Tcp) Listen(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error listen %v, %v", addr, err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error while accept %v", err)
			continue
		}
		session := types.Session{
			Id:          tcp.Id.NewUUID(),
			Conn:        TcpConn{conn: conn},
			ReadBuffer:  types.Buffer{Codec: &types.Json{}},
			WriteBuffer: types.Buffer{Codec: &types.Json{}},
		}
		session.ReadChan = make(chan []byte)
		session.Ticker = time.NewTicker(tcp.HeartBeat * time.Second)
		defer close(session.ReadChan)
		tcp.Sessions[session.Id] = session
		manager.SessionManager().WaitChan <- &session
		log.Printf("clients : %v", tcp.Sessions)
		go tcp.NewClient(&session)
		go session.Read()
	}
}

func (tcp *Tcp) NewClient(session *types.Session) {
	for {
		select {
		case <-session.Ticker.C:
			session.Write(make([]types.Op, 1))
		case msg, ok := <-session.ReadChan:
			if !ok {
				delete(tcp.Sessions, session.Id)
				break
			}
			session.ReadBuffer.Append(msg)
			ops, err := session.ReadBuffer.Package(msg)
			if err != nil {
				log.Printf("Error message : %v", err)
			} else {
				for _, op := range ops {
					op.SetId(session)
					*session.OpPipe <- op
				}
				log.Printf("Recv message : %d", len(ops))
			}
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}
