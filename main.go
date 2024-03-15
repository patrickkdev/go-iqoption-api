package main

import (
	"fmt"
	"math/rand"
	"patrickkdev/Go-IQOption-API/api"
	"patrickkdev/Go-IQOption-API/debug"
	"patrickkdev/Go-IQOption-API/wsapi"
	"time"
)

func main() {
	debug.IfVerbose.Println("Hi mom!")

	email := "patrickfxtrader8q@gmail.com"
	password := "YOUTAP2019"

	userConnection := connectBroker(email, password)
	profile := getProfile(userConnection)

	fmt.Printf("Olá, %s\n", profile.Name)

	startTradingBinaries(userConnection)

	userConnection.WebSocket.WaitGroup.Wait()
}

func startTradingBinaries(userConnection *api.BrokerClient) {
	for {
		tradeDirection := map[bool]wsapi.TradeDirection{
			true:  wsapi.TradeDirectionCall,
			false: wsapi.TradeDirectionPut,
		}[rand.Intn(2) == 0]

		duration := 1

		userConnection.OpenTrade(
			wsapi.TradeTypeBinary,
			100,
			tradeDirection,
			1,
			duration,
			wsapi.TradeBalanceDemo,
		)

		time.Sleep(time.Minute * time.Duration(duration))
	}
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
