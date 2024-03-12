package main

import (
	"fmt"
	"patrickkdev/Go-IQOption-API/api"
	"patrickkdev/Go-IQOption-API/debug"
	"patrickkdev/Go-IQOption-API/wsapi"
)

func main() {
	debug.IfVerbose.Println("Hi mom!")

	userConnection, err := 	api.NewBrokerClient("iqoption.com").
															Login("patrickfxtrader8q@gmail.com", "YOUTAP2019", nil)

	if err != nil {
		panic(err)
	}

	err = userConnection.ConnectSocket()
	if err != nil {
		panic(err)
	}

	coreProfile, err := userConnection.GetCoreProfile()
	if err != nil {
		panic(err)
	}

	println("Login successful")
	fmt.Printf("Hi, %s %s\n", coreProfile.Msg.Result.FirstName, coreProfile.Msg.Result.LastName)

	userProfileClient, err := userConnection.GetProfileClient(int(coreProfile.Msg.Result.ID))
	if err != nil {
		panic(err)
	}

	fmt.Println("Profile: ", userProfileClient.Msg.ImgURL)

	wsapi.TradeDigital(userConnection.WebSocket, 342, 76, 112647980, int(userConnection.TimeSync.GetServerTimestamp()))

	userConnection.WebSocket.WaitGroup.Wait()
}