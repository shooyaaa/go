package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/shooyaaa/codec"
	"github.com/shooyaaa/session"
	"github.com/shooyaaa/uuid"
	"log"
	"net/http"
	"time"
)

type Ws struct {
	Id        uuid.UUID
	Sessions  map[int64]session.Session
	HeartBeat time.Duration
}

func (ws *Ws) Connect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	session := session.Session{
		Id:     ws.Id.NewUUID(),
		Name:   "",
		Conn:   conn,
		Buffer: codec.Buffer{Codec: &codec.Json{}},
	}
	ws.Sessions[session.Id] = session
	log.Printf("clients : %v", ws.Sessions)
	if err != nil {
		log.Print("Upgrade websocket fail :", err)
		return
	}

	session.ReadChan = make(chan []byte)
	defer conn.Close()
	defer close(session.ReadChan)

	conn.SetCloseHandler(session.CloseHandler)
	session.Ticker = time.NewTicker(ws.HeartBeat * time.Second)
	go ws.Read(session)
	for {
		select {
		case <-session.Ticker.C:
			conn.WriteMessage(websocket.PingMessage, []byte{'0', '1', '2', '3'})
		case msg, ok := <-session.ReadChan:
			if !ok {
				delete(ws.Sessions, session.Id)
				break
			}
			req := make(map[string]int)
			session.Buffer.Append(msg)
			err := session.Buffer.Package(req)
			if err != nil {
				log.Printf("Error message : %v", err)
			} else {
				log.Printf("Recv message : %v", req)
			}
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (ws Ws) Read(s session.Session) {
	for {
		if s.Conn == nil {
			break
		}
		_, message, err := s.Conn.(*websocket.Conn).ReadMessage()
		log.Println("read message")
		if err != nil {
			log.Println("Read websocket error :", err)
			break
		} else {
			s.ReadChan <- message
		}
	}
}
