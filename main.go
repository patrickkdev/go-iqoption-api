package main

import (
	"fmt"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/broker"
)

func main() {
	// Mock login data
	loginData := broker.LoginData{
		Email:    "patrickfxtrader8q@gmail.com",
		Password: "YOUTAP2019",
	}

	// Mock host name
	hostName := "iqoption.com"

	// Timeout for requests
	timeout := time.Second * 10
	client := broker.NewClient(loginData, hostName, timeout)

	err := client.Login()
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
	}

	err = client.ConnectSocket()
	if err != nil {
		fmt.Printf("Socket connection failed: %v\n", err)
	}

	for !client.IsReady() {
		time.Sleep(time.Second * 1)
	}

	fmt.Println("Client is ready")

	balances := client.GetBalances()
	if len(balances) == 0 {
		fmt.Println("Failed to retrieve balances")
	} else {
		for _, balance := range balances {
			fmt.Println(balance.Currency, balance.Amount)
		}
	}

	client.OnTradeOpened(func(tradeData broker.TradeData) {
		fmt.Println(tradeData.TimeFrameInMinutes, " minute trade opened: ", tradeData.TradeID)
	})

	newAssets, err := client.GetAssets(broker.AssetTypeDigital)
	if err != nil {
		fmt.Printf("Failed to get assets: %v\n", err)
	} else {
		assets := newAssets.FilterOutNonTradable()

		for i, asset := range assets {
			fmt.Println(i, asset.ActiveID)
		}
	}

	tradeType := broker.AssetTypeDigital
	tradeID := 0
	timeFrame := 1 // Time frame in minutefmt.Println Wait timesync
	time.Sleep(time.Second * 1)

openTrade:
	// Replace parameters with appropriate values
	newTradeID, err := client.OpenTrade(tradeType, 100.0, broker.TradeDirectionCall, 1, timeFrame, broker.BalanceTypeDemo)
	if err != nil {
		fmt.Printf("Failed to open trade: %v\n", err)
	}

	// Ensure trade ID is valid
	if newTradeID == 0 {
		fmt.Println("Invalid trade ID")
	}

	tradeID = newTradeID

	_, err = client.CheckTradeResult(tradeID, timeFrame)
	if err != nil {
		fmt.Printf("Failed to check trade result: %v\n", err)
	}

	if tradeType == broker.AssetTypeDigital {
		tradeType = broker.AssetTypeBinary
		goto openTrade
	}

	for {
		balance, _ := client.GetBalance(broker.BalanceTypeDemo)

		fmt.Printf("Demo Balance: %f\n", balance.Amount)
		time.Sleep(time.Second * 1)
	}
}
