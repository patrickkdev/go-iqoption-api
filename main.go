package main

import (
	"fmt"
	"patrickkdev/Go-IQOption-API/api"
	"patrickkdev/Go-IQOption-API/debug"
	"patrickkdev/Go-IQOption-API/wsapi"
)

func main() {
	debug.IfVerbose.Println("Hi mom!")

	email := "patrickfxtrader8q@gmail.com"
	password := "YOUTAP2019"

	userConnection := connectBroker(email, password)
	profile := getProfile(userConnection)

	fmt.Printf("Olá, %s\n", profile.Name)

	fmt.Println("Profile: ")
	debug.PrintAsJSON(profile)

	balances, err := wsapi.GetBalances(userConnection.WebSocket, int(userConnection.TimeSync.GetServerTimestamp()))
	if err != nil {
		panic(err)
	}

	demoAccBalance, err := balances.FindByType(4)
	if err != nil {
		panic(err)
	}

	debug.PrintAsJSON(demoAccBalance)

	wsapi.TradeBinary(userConnection.WebSocket, 50, 1, int(demoAccBalance.ID), int(userConnection.TimeSync.GetServerTimestamp()))

	userConnection.WebSocket.WaitGroup.Wait()
}

func connectBroker(email string, password string) *api.BrokerClient {
	userConnection, err := api.NewBrokerClient("iqoption.com").Login(email, password, nil)
	if err != nil {
		panic(err)
	}

	err = userConnection.ConnectSocket()
	if err != nil {
		panic(err)
	}

	return userConnection
}

func getProfile(userConnection *api.BrokerClient) *wsapi.UserProfileClient {
	coreProfile, err := userConnection.GetCoreProfile()
	if err != nil {
		panic(err)
	}

	profile, err := userConnection.GetProfileClient(int(coreProfile.Msg.Result.ID))
	if err != nil {
		panic(err)
	}

	return profile
}
