package wsapi

import (
	"fmt"
	"math/rand"
	"patrickkdev/Go-IQOption-API/debug"
	"time"
)

type RequestEvent struct {
	Name      string      `json:"name"`
	Msg       interface{} `json:"msg"`
	RequestId string      `json:"request_id"`
	// Version   string      `json:"version"`
	LocalTime int64       `json:"local_time"`
}

func EmitWithResponse(ws *Socket, event *RequestEvent, responseEventName string, timeout time.Time) (resp []byte, err error) {
	
	event.RequestId = fmt.Sprint(rand.Int63n(10000000000))

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

func CreateUserAvailabilityEvent(activeID int) *RequestEvent {
	msg := map[string]interface{}{
		"name": "update-user-availability",
		"version": "1.1",
		"body": map[string]interface{} {
			"platform_id": "9",
			"idle_duration": 22,
			"selected_asset_id": activeID,
			"selected_asset_type": 7,
		},
	}

	requestEvent := &RequestEvent{
		Name:      "sendMessage",
		Msg:       msg,
		RequestId: "0",
	}

	return requestEvent
}