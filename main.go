package main

import (
	"patrickkdev/Go-IQOption-API/api"
	"patrickkdev/Go-IQOption-API/utils"
	"patrickkdev/Go-IQOption-API/wsapi"
	"time"
)

func main() {
	utils.PrintlnIfVerbose("Hi mom!")

	user, err := api.NewBrokerClient("iqoption.com").
		Login("patrickfxtrader8q@gmail.com", "YOUTAP2019", nil)

	if err != nil {
		panic(err)
	}

	wg, err := user.ConnectSocket()
	if err != nil {
		panic(err)
	}

	user.SendEvent(map[string]interface{}{
		"name": "get-user-profile-client",
	})

	user.Subscribe("heartbeat", func(event wsapi.WSEvent) {
		now := time.Now()
		unixTime := now.Unix()  // Convert to float64 for decimal part
		nanoseconds := int(float64(now.UnixNano()) / 1e9)

		utils.PrintMapAsJSON(event)
		heartbeatAnswer := wsapi.HeartbeatEvent(
			int(event["msg"].(float64)), 
			int(unixTime),
			nanoseconds,
		) 
		utils.PrintMapAsJSON(heartbeatAnswer)

		user.SendEvent(heartbeatAnswer)
	})

	wg.Wait()
}