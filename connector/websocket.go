package connector

import (
	"github.com/shooyaaa/types"
	"github.com/shooyaaa/manager"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type Ws struct {
	Id        		types.UUID
	SessionManager  manager.Session
	HeartBeat 		time.Duration
	Addr	  		string
	Root			string
}


func (ws *Ws) Run() {
	server := HttpServer{
		Root: ws.Root,
		Addr: ws.Addr, 
		Handler: make(map[string]HttpHandler),
	}
	server.Register("/ws", ws.Connect)
	server.Register("/wsinfo", server.Info)
	server.Run()
}

var roomManager = manager.RoomManager{}

func (ws *Ws) Connect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	session := types.Session{
		Id:     ws.Id.NewUUID(),
		Conn:   conn,
		ReadBuffer: types.Buffer{Codec: &types.Json{}},
		WriteBuffer: types.Buffer{Codec: &types.Json{}},
		Status : types.Open,
	}
	if err != nil {
		log.Print("Upgrade websocket fail :", err)
		return
	}
	ws.SessionManager.WaitChan <- session;
	go ws.Read(session)
}

func (ws Ws) Write(session types.Session, i interface{})  error {
	data, _ := session.WriteBuffer.Encode(i)
	return session.Conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, data)
}

func (ws Ws) Broadcast(i interface{}) {
}

func (ws Ws) Read(session types.Session) {
	session.ReadChan = make(chan []byte)
	session.Ticker = time.NewTicker(ws.HeartBeat * time.Millisecond)
	defer session.Ticker.Stop()
	for {
		select {
		case <-session.Ticker.C:
			err := ws.Write(session, []byte{'p', 'i', 'n', 'g'})
			if err != nil {
				log.Printf("Fail to write ping message %d", session.Id)
				session.Status = types.Close
				return
			}
		case msg, ok := <-session.ReadChan:
			if !ok {
				log.Printf("Fail read Readchan : %v", session)
			}
			session.ReadBuffer.Append(msg)
			err := session.ReadBuffer.Package(session.OpPipe, msg)
			if err != nil {
				log.Printf("Error message : %v", err)
			} else {
				log.Printf("Recv message : %v", msg)
			}
		default:
			buffer := make([]byte, 1024)
			_, buffer, err := session.Conn.(*websocket.Conn).ReadMessage()
			if err != nil {
				log.Printf("Error while Read msg %v", err)
			} else {
				session.ReadChan <- buffer
			}
			time.Sleep(5 * time.Millisecond)
		}
	}
}
