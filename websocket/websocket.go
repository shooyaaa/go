package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/shooyaaa/uuid"
	"log"
	"net/http"
)

type Ws struct {
	Addr string
}

var id = uuid.Simple{}
var clients = make(map[int64]Session)

func (ws *Ws) Connect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	session := Session{
		id.NewUUID(),
		"",
		conn,
	}
	clients[session.Id] = session
	log.Printf("clients : %v", clients)
	if err != nil {
		log.Print("Upgrade websocket fail :", err)
		return
	}
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read websocket error :", err)
			break
		}
		log.Printf("Recv message : %s", message)
	}
}
