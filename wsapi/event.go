package wsapi

import (
	"fmt"
	"patrickkdev/Go-IQOption-API/debug"
	"time"
)

type Event struct {
	Name      string      `json:"name"`
	Msg       interface{} `json:"msg"`
	Version   string      `json:"version"`
	RequestId string      `json:"request_id"`
	LocalTime int64       `json:"local_time"`
}

func NewEvent(name string, msg map[string]interface{}, requestId string) *Event {
	return &Event{
		Name:      name,
		Msg:       msg,
		Version:   "1.0",
		RequestId: requestId,
	}
}

func EmitWithResponse(ws *Socket, event *Event, responseEventName string, timeout time.Time) (resp []byte, err error) {
	ws.EmitEvent(event)

	ws.AddEventListener(responseEventName, func(responseEvent []byte) {
		resp = responseEvent
	})

	debug.IfVerbose.Println("waiting for " + responseEventName + "...")

	for {
		if resp != nil {
			debug.IfVerbose.Println("response " + responseEventName + " found")
			break
		}

		if time.Since(timeout) > 0 {
			err = fmt.Errorf("timed out waiting for response" + responseEventName)
			debug.IfVerbose.Println(err.Error())
			break
		}
	}

	ws.RemoveEventListener(responseEventName)
	debug.IfVerbose.PrintAsJSON(resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
