package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/api"
	"github.com/patrickkdev/Go-IQOption-API/data"
	"github.com/patrickkdev/Go-IQOption-API/debug"
	"github.com/patrickkdev/Go-IQOption-API/wsapi"
)

func main() {
	debug.IfVerbose.Println("Hi mom!")

	email := "patrickfxtrader8q@gmail.com"
	password := "YOUTAP2019"

	userConnection := connectBroker(email, password)
	profile := getProfile(userConnection)

	fmt.Printf("Olá, %s\n", profile.Msg.UserName)

	startTradingBinaries(userConnection)

	userConnection.WS.WaitGroup.Wait()
}

func startTradingBinaries(userConnection *api.BrokerClient) {
	for {
		if userConnection.WS.Closed {
			fmt.Println("No connection with socket")

			time.Sleep(time.Second)
			continue
		}

		tradeType := map[bool]wsapi.TradeType{
			true:  wsapi.TradeTypeDigital,
			false: wsapi.TradeTypeBinary,
		}[rand.Float32() < 0.5]

		duration := rand.Intn(4) + 1

		pair, err := data.Pairs.GetByName("EURUSD-OTC")
		if err != nil {
			fmt.Println(err)
			continue
		}

		candles, err := userConnection.GetCandles(20, 1, time.Now().UnixMicro(), pair)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Color strategy, trade call if last candle was green
		// and put if last candle was red
		tradeDirection := wsapi.TradeDirectionCall
		if candles[0].Close < candles[0].Open {
			tradeDirection = wsapi.TradeDirectionPut
		}

		fmt.Printf("Opening %s trade. Direction: %s. Pair: %d. Duration: %d minutes\n", tradeType, tradeDirection, pair, duration)

		_, win, err := userConnection.OpenTrade(
			tradeType,
			2,
			tradeDirection,
			pair,
			duration,
			wsapi.TradeBalanceDemo,
			true,
		)

		if err != nil {
			fmt.Println("Error opening trade:", err)
			continue
		}

		tradeResult := map[bool]string{
			true:  "win",
			false: "lose",
		}[win]

		fmt.Printf("Trade result: %s\n", tradeResult)

		time.Sleep(time.Second)
	}
}

func connectBroker(email string, password string) *api.BrokerClient {
	userConnection, err := api.NewBrokerClient("iqoption.com", time.Duration(time.Second*10)).
		Login(email, password, nil)

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
