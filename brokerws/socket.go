package brokerws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/internal/debug"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type EventCallback func(event []byte)

type Socket struct {
	Conn          *websocket.Conn
	Closed        bool
	eventHandlers map[string]EventCallback
	WaitGroup     *sync.WaitGroup
}

const timeout = time.Second * 15

func NewSocketConnection(url string, onLoseConnection func()) (*Socket, error) {
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
		debug.IfVerbose.Println("conn is nil")
		return nil, fmt.Errorf("conn is nil")
	}

	conn.SetReadLimit(-1)

	newSocketConnection := &Socket{
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

func (ws *Socket) Close() {
	ws.Conn.Close(websocket.StatusNormalClosure, "close")
}

func (ws *Socket) EmitEvent(event interface{}) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*15)

	wsjson.Write(ctx, ws.Conn, event)
	ctxCancel()
}

func (ws *Socket) AddEventListener(name string, callback EventCallback) {
	ws.eventHandlers[name] = callback
}

func (ws *Socket) RemoveEventListener(name string) {
	delete(ws.eventHandlers, name)
}

func (ws *Socket) handleEvent(eventB []byte) {
	reportEventError := func(errMessage string) {
		debug.IfVerbose.Println(errMessage)
		debug.IfVerbose.PrintAsJSON(eventB)
	}

	event := new(struct {
		Name   string      `json:"name"`
		Result interface{} `json:"result"`
	})

	err := json.Unmarshal(eventB, &event)

	if err != nil {
		reportEventError("error unmarshalling event")
		return
	}

	callback, ok := ws.eventHandlers[event.Name]
	if !ok {
		reportEventError("no callback found for event: " + event.Name)
		return
	}

	callback(eventB)
}

func (ws *Socket) Listen(onLoseConnection func()) {
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
