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

type WsConn struct {
	conn *websocket.Conn
}

func (wc WsConn) Read() ([]byte, error) {
	_, buffer, err := wc.conn.ReadMessage()
	return buffer, err
}

func (wc WsConn) Write(bytes []byte) (int, error) {
	err := wc.conn.WriteMessage(websocket.TextMessage, bytes)
	count := 0
	if err == nil {
		count = len(bytes)
	}
	return count, err
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
		Conn:        WsConn{conn: conn},
		ReadBuffer:  types.Buffer{Codec: &types.Json{}},
		WriteBuffer: types.Buffer{Codec: &types.Json{}},
		Status:      types.Waiting,
	}
	if err != nil {
		log.Print("Upgrade websocket fail :", err)
		return
	}
	log.Printf("Incoming Session %d", session.Id)
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
			_, err := session.Write(ops)
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
			if session.Status == types.Pending {
				session.Status = types.Open
				go session.Read()
			} else if session.Status == types.Close {
				log.Printf("Session closed, so stop communicate with it")
				return
			}
			time.Sleep(5 * time.Millisecond)
		}
	}
}
