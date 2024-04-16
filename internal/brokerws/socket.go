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

func NewSocketConnection(url string, timeoutDuration time.Duration, onLoseConnection func()) (*WebSocket, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeoutDuration)
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

	go newSocketConnection.Listen(timeoutDuration, onLoseConnection)

	return newSocketConnection, nil
}

func (ws *WebSocket) IsConnectionOK() bool {
	return !ws.Closed
}

func (ws *WebSocket) Close() {
	ws.Conn.Close(websocket.StatusNormalClosure, "close")
	ws.Closed = true
}

func (ws *WebSocket) Listen(timeoutDuration time.Duration, onLoseConnection func()) {
	const maxErrorCount = 5
	var errorCount int

	for ws.IsConnectionOK() {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()

		if errorCount > maxErrorCount {
			debug.IfVerbose.Println("too many errors, closing ws connection")
			ws.Conn.Close(websocket.StatusAbnormalClosure, "close")
			ws.Closed = true

			onLoseConnection()
			break
		}

		_, message, err := ws.Conn.Read(ctx)
		cancel()

		if err != nil {
			if websocket.CloseStatus(err) != websocket.StatusNormalClosure {
				onLoseConnection()

				// Connection closed by remote host
				break
			}

			debug.IfVerbose.Println("error reading ws event:", err)
			errorCount++
			continue
		}

		// Reset error count on successful read
		errorCount = 0

		ws.handleEvent(message)
	}
}
