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

	userConnection, err := api.NewBrokerClient("iqoption.com").
		Login("patrickfxtrader8q@gmail.com", "YOUTAP2019", nil)

	if err != nil {
		panic(err)
	}

	wg, err := userConnection.ConnectSocket()
	if err != nil {
		panic(err)
	}

	userConnection.Subscribe("heartbeat", func(event wsapi.Event) {
		now := time.Now()
		unixTime := now.UnixNano()
		requestId := fmt.Sprint(unixTime)[10:18]

		heartbeatTime := int(event.Msg.(float64))

		utils.PrintlnIfVerbose("Received heartbeat event at:", int(event.Msg.(float64)))

		heartbeatAnswer := wsapi.NewEvent(
			"heartbeat",
			map[string]interface{}{
				"heartbeatTime": heartbeatTime,
				"userTime":      int(userConnection.TimeSync.GetServerTimestamp() * 1000),
			},
			requestId,
		)

		userConnection.SendEvent(heartbeatAnswer)
		utils.PrintMapAsJSON(heartbeatAnswer.Msg.(map[string]interface{}))
	})

	wg.Wait()
}
