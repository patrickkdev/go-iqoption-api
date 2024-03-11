package wsapi

import "fmt"

func HeartbeatEvent(heartBeatTime int, userTime int, requestId string) (hbEvent map[string]interface{}) {
	event := map[string]interface{}{
		"name": "heartbeat",
		"msg": map[string]string{
			"heartbeatTime": fmt.Sprint(heartBeatTime),
			"userTime":      fmt.Sprint(userTime),
		},
		"request_id": requestId,
	}
	return event
}
