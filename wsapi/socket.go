package wsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"patrickkdev/Go-IQOption-API/debug"
	"sync"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type EventCallback func(event []byte)

type Socket struct {
	Conn          *websocket.Conn
	Closed        bool
	eventHandlers map[string]EventCallback
	ctx           context.Context
	WaitGroup     *sync.WaitGroup
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
		debug.IfVerbose.Println("conn is nil")
		return nil, fmt.Errorf("conn is nil")
	}

	conn.SetReadLimit(-1)

	newSocketConnection := &Socket{
		eventHandlers: make(map[string]EventCallback),
		ctx:           ctx,
		Conn:          conn,
		Closed:        false,
		WaitGroup:     new(sync.WaitGroup),
	}

	debug.IfVerbose.Println("new socket connection: ")

	newSocketConnection.WaitGroup.Add(1)

	go newSocketConnection.Listen()

	return newSocketConnection, nil
}

func (ws *Socket) Close() {
	ws.Conn.Close(websocket.StatusNormalClosure, "close")
}

func (ws *Socket) EmitEvent(event interface{}) {
	wsjson.Write(ws.ctx, ws.Conn, event)
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

func (ws *Socket) Listen() {
	defer ws.WaitGroup.Done()

	var errorCount int = 0

	for {
		if ws.Closed {
			break
		}

		if errorCount > 5 {
			debug.IfVerbose.Println("too many errors")
			ws.Conn.Close(websocket.StatusAbnormalClosure, "close")
			ws.Closed = true
			continue
		}

		_, message, err := ws.Conn.Read(ws.ctx)

		if err != nil {
			debug.IfVerbose.Println("error reading ws event:", err)
			errorCount++
			continue
		}

		ws.handleEvent(message)
	}
}
