package brokerws

import (
	"encoding/json"

	"github.com/patrickkdev/Go-IQOption-API/btypes"
)

func NewTimesync() btypes.Timesync {
	return btypes.Timesync{
		Name:            "timeSync",
		ServerTimestamp: 0,
		ExpirationTime:  0,
	}
}

func (ws *Socket) OnTimeSync(callback func(newTimestamp int64)) {
	ws.AddEventListener("timeSync", func(event []byte) {
		var timesyncEvent btypes.TimesyncEvent
		json.Unmarshal(event, &timesyncEvent)

		callback(timesyncEvent.Msg)
	})
}
