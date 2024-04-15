package brokerws

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/btypes"
	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
)

func (ws *Socket) EmitWithResponse(event *btypes.RequestEvent, responseEventName string, timeout time.Time) (resp []byte, err error) {
	event.RequestId = fmt.Sprint(rand.Int63n(10000000000))

	ws.EmitEvent(event)

	ws.AddEventListener(responseEventName, func(responseEvent []byte) {
		resp = responseEvent
	})

	debug.IfVerbose.Println("waiting for " + responseEventName + "...")

	for {
		if ws.Closed {
			err = fmt.Errorf("websocket connection closed :: func EmitWithResponse (" + responseEventName + ")")
			break
		}

		if time.Since(timeout) > 0 {
			err = fmt.Errorf("timed out waiting for response: " + responseEventName)
			debug.IfVerbose.Println(err.Error())
			break
		}

		if resp != nil {
			debug.IfVerbose.Println("response " + responseEventName + " found")
			break
		}
	}

	ws.RemoveEventListener(responseEventName)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ws *Socket) Emit(event *btypes.RequestEvent) {
	event.RequestId = fmt.Sprint(rand.Int63n(10000000000))

	ws.EmitEvent(event)
}
