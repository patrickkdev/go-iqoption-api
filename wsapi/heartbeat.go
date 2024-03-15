package wsapi

import (
	"patrickkdev/Go-IQOption-API/debug"
)

type Heartbeat struct {
	Name string `json:"name"`
	Msg  int 		`json:"msg"`
}

func AnswerHeartbeat(ws *Socket, heartbeatFromServer Heartbeat, serverTimestamp int64) {	
		heartbeatTime := int(heartbeatFromServer.Msg)

		debug.IfVerbose.Println("Heartbeat:", int(heartbeatTime))
	
		heartbeatFromClient := &RequestEvent{
			Name:      "heartbeat",
			Msg:       map[string]interface{}{
				"heartbeatTime": heartbeatTime,
				"userTime":      serverTimestamp * 1000,
			},
		}
	
		Emit(ws, heartbeatFromClient)
}
