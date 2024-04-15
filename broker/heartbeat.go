package broker

import (
	"encoding/json"

	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
	"github.com/patrickkdev/Go-IQOption-API/internal/types"
)

func (c *Client) startAnsweringHeartBeats() {
	debug.IfVerbose.Println("Answering heartbeats...")

	c.ws.AddEventListener("heartbeat", func(event []byte) {
		var heartbeat types.Heartbeat
		json.Unmarshal(event, &heartbeat)

		heartbeatAnswer := requestEvent{
			Name: "heartbeat",
			Msg: map[string]interface{}{
				"heartbeatTime": int(heartbeat.Msg),
				"userTime":      c.serverTimestamp * 1000,
			},
		}.WithRandomRequestId()

		c.ws.Emit(heartbeatAnswer)
	})
}
