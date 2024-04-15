package brokerws

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/internal/debug"

	"nhooyr.io/websocket"
)

type EventCallback func(event []byte)

type WebSocket struct {
	Conn               *websocket.Conn
	Closed             bool
	eventHandlers      map[string]EventCallback
	eventHandlersMutex sync.RWMutex
	WaitGroup          *sync.WaitGroup
}

const timeout = time.Second * 15

func NewSocketConnection(url string, onLoseConnection func()) (*WebSocket, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	defer ctxCancel()

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
		err = fmt.Errorf("conn is nil")
		debug.IfVerbose.Println(err.Error())
		return nil, err
	}

	conn.SetReadLimit(-1)

	newSocketConnection := &WebSocket{
		eventHandlers: make(map[string]EventCallback),
		Conn:          conn,
		Closed:        false,
		WaitGroup:     new(sync.WaitGroup),
	}

	debug.IfVerbose.Println("new socket connection: ")

	newSocketConnection.WaitGroup.Add(1)

	go newSocketConnection.Listen(onLoseConnection)

	return newSocketConnection, nil
}

func (ws *WebSocket) IsConnectionOK() bool {
	return !ws.Closed
}

func (ws *WebSocket) Close() {
	ws.Conn.Close(websocket.StatusNormalClosure, "close")
}

func (ws *WebSocket) Listen(onLoseConnection func()) {
	var errorCount int = 0

	for {
		ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
		defer ctxCancel()

		if errorCount > 5 {
			debug.IfVerbose.Println("too many errors")
			ws.Conn.Close(websocket.StatusAbnormalClosure, "close")
			ws.Closed = true
		}

		if ws.Closed {
			onLoseConnection()
			break
		}

		_, message, err := ws.Conn.Read(ctx)
		ctxCancel()

		if err != nil {
			debug.IfVerbose.Println("error reading ws event:", err)
			errorCount++
			continue
		}

		ws.handleEvent(message)
	}
}
