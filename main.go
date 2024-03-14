package main

import (
	"fmt"
	"math/rand"
	"patrickkdev/Go-IQOption-API/api"
	"patrickkdev/Go-IQOption-API/data"
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

	fmt.Printf("Olá, %s\n", profile.Msg.UserName)

	startTradingBinaries(userConnection)

	userConnection.WebSocket.WaitGroup.Wait()
}

func startTradingBinaries(userConnection *api.BrokerClient) {
	for {
		tradeType := map[bool]wsapi.TradeType{
			true:  wsapi.TradeTypeDigital,
			false: wsapi.TradeTypeBinary,
		}[rand.Float32() < 0.5]

		duration := 1// rand.Intn(4) + 1
		
		pair, err := data.Pairs.GetByName("EURUSD-OTC")
		if err != nil {
			panic(err)
		}

		candles, err := wsapi.GetCandles(userConnection.WebSocket, 2, 60 * 1, int(time.Now().UnixMicro()), pair)
		if err != nil {
			panic(err)
		}

		// Color strategy, trade call if last candle was green
		// and put if last candle was red
		tradeDirection := wsapi.TradeDirectionCall
		if candles[0].Close < candles[0].Open {
			tradeDirection = wsapi.TradeDirectionPut
		}

		tradeID, err := userConnection.OpenTrade(
			tradeType,
			2,
			tradeDirection,
			pair,
			duration,
			wsapi.TradeBalanceDemo,
		)
		if err != nil {
			fmt.Println(err)
			continue
		}
		
		fmt.Printf("Successully opened trade of type '%s' with ID '%d'\n", tradeType, tradeID)

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
