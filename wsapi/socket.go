package wsapi

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type EventCallback func(event Event)

type Socket struct {
	Conn   *websocket.Conn
	Closed bool
	ctx    context.Context
}

func NewSocketConnection(url string) (*Socket, error) {
	ctx := context.Background()

	var header http.Header = http.Header{}
	header.Add("check_hostname", "false")
	header.Add("cert_reqs", "none")
	header.Add("ca_certs", "cacert.pem")

	conn, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		HTTPHeader: header,
	})
	if err != nil {
		return nil, err
	}

	if conn == nil {
		println("conn is nil")
	}

	newSocketConnection := &Socket{
		Conn:   conn,
		Closed: false,
		ctx:    ctx,
	}

	println("new socket connection: ")

	return newSocketConnection, nil
}

func (ws *Socket) Close() {
	ws.Conn.Close(websocket.StatusNormalClosure, "close")
}

func (ws *Socket) Write(event interface{}) {
	wsjson.Write(ws.ctx, ws.Conn, event)
}

func (ws *Socket) Listen(wg *sync.WaitGroup, handleEvent EventCallback) {
	defer wg.Done()

	var errorCount int = 0
	var message Event

	for {
		if ws.Closed {
			break
		}

		if errorCount > 5 {
			println("too many errors")
			ws.Conn.Close(websocket.StatusAbnormalClosure, "close")
			ws.Closed = true
			continue
		}

		err := wsjson.Read(ws.ctx, ws.Conn, &message)
		if err != nil {
			fmt.Println("error reading ws event:", err)
			errorCount++
			continue
		}

		handleEvent(message)
	}
}
