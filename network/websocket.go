package network

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shooyaaa/types"
)

type Ws struct {
	Id        types.UUID
	HeartBeat time.Duration
	Root      string
	waitChan  chan *types.Session
	server    HttpServer
}

type WsConn struct {
	conn *websocket.Conn
}

func (wc WsConn) Read(buffer []byte) (int, error) {
	count, buffer, err := wc.conn.ReadMessage()
	return count, err
}

func (wc WsConn) Write(bytes []byte) (int, error) {
	err := wc.conn.WriteMessage(websocket.TextMessage, bytes)
	count := 0
	if err == nil {
		count = len(bytes)
	}
	return count, err
}

func (ws *Ws) Listen(addr string) error {
	ws.server = HttpServer{
		Root:    ws.Root,
		Addr:    addr,
		Handler: make(map[string]HttpHandler),
	}
	ws.server.Register("/ws", ws.Connect)
	ws.server.Register("/wsinfo", ws.server.Info)
	ws.server.Run()
	return nil
}

func (ws *Ws) Connect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	session := types.Session{
		Id:   ws.Id.NewUUID(),
		Conn: WsConn{conn: conn},
	}
	if err != nil {
		log.Print("Upgrade websocket fail :", err)
		return
	}
	log.Printf("Incoming Session %d", session.Id)
	ws.waitChan <- &session
	//go ws.Commuicate(&session)
}

func (ws Ws) Accept() *types.Session {
	return <-ws.waitChan
}
func (ws Ws) Close() error {
	ws.server.Close()
	return nil
}

/*
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


*/
