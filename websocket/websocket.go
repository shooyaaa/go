package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/shooyaaa/uuid"
	"log"
	"net/http"
	"time"
)

type Ws struct {
	Id        uuid.UUID
	Sessions  map[int64]Session
	HeartBeat time.Duration
}

func (ws *Ws) Connect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	session := Session{
		Id:   ws.Id.NewUUID(),
		Name: "",
		Conn: conn,
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

	conn.SetCloseHandler(session.closeHandler)
	session.Ticker = time.NewTicker(ws.HeartBeat * time.Second)
	go session.Read(session.ReadChan)
	for {
		select {
		case <-session.Ticker.C:
			log.Println("start send ping message")
			conn.WriteMessage(websocket.PingMessage, []byte{'0', '1', '2', '3'})
		case msg, ok := <-session.ReadChan:
			if !ok {
				delete(ws.Sessions, session.Id)
				break
			}
			log.Printf("Recv message : %s", msg)
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}
