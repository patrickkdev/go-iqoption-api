package wsapi

import (
	"context"
	"fmt"
	"sync"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WSEvent map[string]interface{}

type WSEventCallback func(event WSEvent)

type WSocket struct {
	Conn *websocket.Conn
	Closed bool
	ctx context.Context
}

func NewSocketConnection(url string) (*WSocket, error) {
	ctx := context.Background()

	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	if conn == nil {
		println("conn is nil")
	}

	newSocketConnection := &WSocket{
		Conn: conn,
		Closed: false,
		ctx: ctx,
	}

	println("new socket connection: ")

	return newSocketConnection, nil
}

func (ws *WSocket) Close() {
	ws.Conn.Close(websocket.StatusNormalClosure, "close")
}

func (ws *WSocket) Write(event interface{}) () {
	wsjson.Write(ws.ctx, ws.Conn, event)
}

func (ws *WSocket) Listen(wg *sync.WaitGroup, handleEvent WSEventCallback) () {
	defer wg.Done()

	var errorCount int = 0
	var message map[string]interface{}

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