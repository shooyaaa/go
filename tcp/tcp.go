package tcp

import (
	"github.com/shooyaaa/codec"
	"github.com/shooyaaa/session"
	"github.com/shooyaaa/uuid"
	"log"
	"net"
	"time"
)

type Tcp struct {
	Id        uuid.UUID
	Sessions  map[int64]session.Session
	HeartBeat time.Duration
}

func (tcp *Tcp) Listen(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Error listen %v, %v", addr, err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error while accept %v", err)
			continue
		}
		session := session.Session{
			Id:     tcp.Id.NewUUID(),
			Name:   "",
			Conn:   conn,
			Buffer: codec.Buffer{Codec: &codec.Json{}},
		}
		session.ReadChan = make(chan []byte)
		session.Ticker = time.NewTicker(tcp.HeartBeat * time.Second)
		defer close(session.ReadChan)
		tcp.Sessions[session.Id] = session
		log.Printf("clients : %v", tcp.Sessions)
		go tcp.NewClient(session)
		go tcp.Read(session)
	}
}

func (tcp *Tcp) Read(s session.Session) {
	for {
		if s.Conn == nil {
			break
		}
		data := make([]byte, 100)
		_, err := s.Conn.(*net.TCPConn).Read(data)
		if err != nil {
			log.Println("Read tcp error :", err)
			break
		} else {
			s.ReadChan <- data
		}
	}

}

func (tcp *Tcp) NewClient(session session.Session) {
	for {
		select {
		case <-session.Ticker.C:
			session.Conn.(*net.TCPConn).Write([]byte{'0', '1', '2', '3'})
		case msg, ok := <-session.ReadChan:
			if !ok {
				delete(tcp.Sessions, session.Id)
				break
			}
			req := make(map[string]interface{})
			session.Buffer.Append(msg)
			err := session.Buffer.Package(req)
			if err != nil {
				log.Printf("Error message : %v : %v", len(msg), err)
			} else {
				log.Printf("Recv message : %s", req)
			}
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}
