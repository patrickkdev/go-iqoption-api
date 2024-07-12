package brokerws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
	"nhooyr.io/websocket/wsjson"
)

func (ws *WebSocket) Subscribe(event any, responseEventName string, callback EventCallback) {
	ws.Emit(event)

	ws.AddEventListener(responseEventName, callback)
}

func (ws *WebSocket) Emit(event any) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*15)
	defer ctxCancel()

	wsjson.Write(ctx, ws.Conn, event)
}

func (ws *WebSocket) EmitWithResponse(event any, responseEventName string, timeout time.Duration) (resp []byte, err error) {
	debug.IfVerbose.Println("Emiting event: ")
	debug.IfVerbose.PrintAsJSON(event)

	ws.Emit(event)

	ch := make(chan []byte)
	ws.AddEventListener(responseEventName, func(responseEvent []byte) {
		ch <- responseEvent
	})
	debug.IfVerbose.Println("waiting for " + responseEventName + "...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case resp = <-ch:
		debug.IfVerbose.Println("Received '" + responseEventName + "' event")
	case <-ctx.Done():
		err = fmt.Errorf("timed out waiting for response: %s", responseEventName)
		debug.IfVerbose.Println(err.Error())
	}

	ws.RemoveEventListener(responseEventName)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ws *WebSocket) AddEventListener(name string, callback EventCallback) {
	ws.eventHandlers.Store(name, callback)
}

func (ws *WebSocket) RemoveEventListener(name string) {
	ws.eventHandlers.Delete(name)
}

func (ws *WebSocket) handleEvent(eventB []byte) {
	reportEventError := func(errMessage string) {
		debug.IfVerbose.Println(errMessage)
	}

	event := new(struct {
		Name string `json:"name"`
	})

	err := json.Unmarshal(eventB, &event)

	// Ignore heartbeat events
	if event.Name == "heartbeat" {
		return
	}

	if err != nil {
		reportEventError("error unmarshalling event")
		return
	}

	value, ok := ws.eventHandlers.Load(event.Name)
	if !ok {
		reportEventError("no callback found for event: " + event.Name)
		return
	}

	callback, ok := value.(EventCallback)
	if !ok {
		reportEventError("invalid callback type for event: " + event.Name)
		return
	}

	callback(eventB)
}
