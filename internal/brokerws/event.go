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

	wsjson.Write(ctx, ws.Conn, event)
	ctxCancel()
}

func (ws *WebSocket) EmitWithResponse(event any, responseEventName string, timeout time.Time) (resp []byte, err error) {
	ws.Emit(event)

	ws.AddEventListener(responseEventName, func(responseEvent []byte) {
		resp = responseEvent
	})

	debug.IfVerbose.Println("waiting for " + responseEventName + "...")

	for {
		if !ws.IsConnectionOK() {
			err = fmt.Errorf("websocket connection closed :: func EmitWithResponse (" + responseEventName + ")")
			break
		}

		if time.Now().After(timeout) {
			err = fmt.Errorf("timed out waiting for response: " + responseEventName)
			debug.IfVerbose.Println(err.Error())
			break
		}

		if resp != nil {
			debug.IfVerbose.Println("Received '" + responseEventName + "' event")
			break
		}
	}

	ws.RemoveEventListener(responseEventName)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ws *WebSocket) AddEventListener(name string, callback EventCallback) {
	ws.eventHandlersMutex.Lock()
	ws.eventHandlers[name] = callback
	ws.eventHandlersMutex.Unlock()
}

func (ws *WebSocket) RemoveEventListener(name string) {
	ws.eventHandlersMutex.Lock()
	delete(ws.eventHandlers, name)
	ws.eventHandlersMutex.Unlock()
}

func (ws *WebSocket) handleEvent(eventB []byte) {
	reportEventError := func(errMessage string) {
		debug.IfVerbose.Println(errMessage)
	}

	event := new(struct {
		Name   string `json:"name"`
		Result any    `json:"result"`
	})

	err := json.Unmarshal(eventB, &event)

	if err != nil {
		reportEventError("error unmarshalling event")
		return
	}

	ws.eventHandlersMutex.Lock()
	callback, ok := ws.eventHandlers[event.Name]
	ws.eventHandlersMutex.Unlock()

	if !ok {
		reportEventError("no callback found for event: " + event.Name)
		return
	}

	callback(eventB)
}
