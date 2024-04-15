package brokerws

import (
	"encoding/json"

	"github.com/patrickkdev/Go-IQOption-API/btypes"
)

func (ws *Socket) OnHeartBeat(callback func(heartbeat btypes.Heartbeat)) {
	ws.AddEventListener("heartbeat", func(event []byte) {
		var heartbeat btypes.Heartbeat
		json.Unmarshal(event, &heartbeat)

		callback(heartbeat)
	})
}

func (ws *Socket) AnswerHeartBeat(heartbeatFromServer btypes.Heartbeat, serverTimestamp int64) {
	heartbeatAnswer := &btypes.RequestEvent{
		Name: "heartbeat",
		Msg: map[string]interface{}{
			"heartbeatTime": int(heartbeatFromServer.Msg),
			"userTime":      serverTimestamp * 1000,
		},
	}

	ws.Emit(heartbeatAnswer)
}
