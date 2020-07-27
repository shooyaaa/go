package connector

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shooyaaa/manager"
	"github.com/shooyaaa/types"
)

type Ws struct {
	Id        types.UUID
	HeartBeat time.Duration
	Addr      string
	Root      string
}

func (ws *Ws) Run() {
	server := HttpServer{
		Root:    ws.Root,
		Addr:    ws.Addr,
		Handler: make(map[string]HttpHandler),
	}
	server.Register("/ws", ws.Connect)
	server.Register("/wsinfo", server.Info)
	server.Run()
}

func (ws *Ws) Connect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	session := types.Session{
		Id:          ws.Id.NewUUID(),
		Conn:        conn,
		ReadBuffer:  types.Buffer{Codec: &types.Json{}},
		WriteBuffer: types.Buffer{Codec: &types.Json{}},
		Status:      types.Waiting,
	}
	if err != nil {
		log.Print("Upgrade websocket fail :", err)
		return
	}
	log.Print("Incoming Session %d", session.Id)
	manager.SessionManager().WaitChan <- &session
	go ws.Commuicate(&session)
}

func (ws Ws) Broadcast(i interface{}) {
}

func (ws Ws) Commuicate(session *types.Session) {
	session.ReadChan = make(chan []byte)
	session.Ticker = time.NewTicker(ws.HeartBeat * time.Millisecond)
	defer session.Ticker.Stop()
	for {
		select {
		case <-session.Ticker.C:
			ops := make([]types.Op, 1)
			err := session.Write(ops)
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
			ops, err := session.ReadBuffer.Package(msg)
			for _, op := range ops {
				op.SetId(session)
				*session.OpPipe <- op
			}
			if err != nil {
				log.Printf("Error message : %v", err)
			} else {
				log.Printf("Recv message : %v", msg)
			}
		default:
			if session.Status == types.Pending {
				session.Status = types.Open
				go ws.Read(session)
			}
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func (ws Ws) Read(session *types.Session) {
	for {
		_, buffer, err := session.Conn.(*websocket.Conn).ReadMessage()
		if err != nil {
			log.Printf("Error while Read msg %v", err)
			break
		} else {
			session.ReadChan <- buffer
		}
	}
}
