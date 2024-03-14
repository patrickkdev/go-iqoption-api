package wsapi

import (
	"fmt"
	"patrickkdev/Go-IQOption-API/debug"
	"time"
)

type Heartbeat struct {
	Name string `json:"name"`
	Msg  int 		`json:"msg"`
}

func AnswerHeartbeat(ws *Socket, heartbeatFromServer Heartbeat, serverTimestamp int64) {
		now := time.Now()
		unixTime := now.UnixNano()
		requestId := fmt.Sprint(unixTime)[10:18]
	
		heartbeatTime := int(heartbeatFromServer.Msg)

		debug.If(bool(debug.IfVerbose) && false).Println("Received heartbeat from server at:", int(heartbeatTime))
	
		heartbeatFromClient := &RequestEvent{
			RequestId: requestId,
			Name:      "heartbeat",
			Msg:       map[string]interface{}{
				"heartbeatTime": heartbeatTime,
				"userTime":      serverTimestamp * 1000,
			},
		}
	
		ws.EmitEvent(heartbeatFromClient)
}
