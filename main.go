package main

import (
	"fmt"
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
		unixTime := now.UnixNano() // Convert to float64 for decimal part
		requestId := fmt.Sprint(unixTime)[10:18]

		utils.PrintMapAsJSON(event)
		heartbeatAnswer := wsapi.HeartbeatEvent(
			int(event["msg"].(float64)),
			int(unixTime),
			requestId,
		)
		utils.PrintMapAsJSON(heartbeatAnswer)

		user.SendEvent(heartbeatAnswer)
	})

	wg.Wait()
}
