package broker

import (
	"encoding/json"

	"github.com/patrickkdev/Go-IQOption-API/internal/types"
)

func (c *Client) keepServerTimestampUpdated() {
	c.ws.AddEventListener("timeSync", func(event []byte) {
		var timesyncEvent types.TimesyncEvent
		json.Unmarshal(event, &timesyncEvent)

		c.serverTimestamp = timesyncEvent.Msg / 1000
	})
}
