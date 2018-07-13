package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/shooyaaa/arena"
	"github.com/shooyaaa/types"
	"github.com/shooyaaa/uuid"
	"log"
	"net/http"
	"time"
)

type Ws struct {
	Id        uuid.UUID
	Sessions  map[uuid.ID]types.Session
	HeartBeat time.Duration
}

var roomManager = arena.RoomManager{}

func (ws *Ws) Connect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	session := types.Session{
		Id:     ws.Id.NewUUID(),
		Name:   "",
		Conn:   conn,
		Buffer: types.Buffer{Codec: &types.Json{}},
	}
	ws.Sessions[session.Id] = session
	log.Printf("clients : %v", ws.Sessions)
	if err != nil {
		log.Print("Upgrade websocket fail :", err)
		return
	}

	session.ReadChan = make(chan []byte)

	session.Ticker = time.NewTicker(ws.HeartBeat * time.Second)
	defer session.Ticker.Stop()
	go ws.Read(session)
	for {
		select {
		case <-session.Ticker.C:
			now := time.Now().UnixNano()
			ws.Write(session, types.OpPing{"ping", now})
		case msg, ok := <-session.ReadChan:
			if !ok {
				log.Printf("Fail read Readchan : %v", ws.Sessions)
				delete(ws.Sessions, session.Id)
				log.Printf("Fail read Readchan : %v", ws.Sessions)
				return
			}
			session.Buffer.Append(msg)
			err := session.Buffer.Package(types.Dispatcher(roomManager), msg)
			if err != nil {
				log.Printf("Error message : %v", err)
			} else {
				log.Printf("Recv message : %v", msg)
			}
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (ws Ws) Write(session types.Session, i interface{}) {
	data, _ := session.Buffer.Encode(i)
	session.Conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, data)
}

func (ws Ws) Broadcast(i interface{}) {
	for _, session := range ws.Sessions {
		ws.Write(session, i)
	}
}

func (ws Ws) Read(s types.Session) {
	for {
		if s.Conn == nil {
			break
		}
		_, message, err := s.Conn.(*websocket.Conn).ReadMessage()
		log.Println("read message")
		if err != nil {
			close(s.ReadChan)
			log.Println("Read websocket error :", err)
			break
		} else {
			s.ReadChan <- message
		}
	}
}
