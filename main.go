package main

import (
	"fmt"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/broker"
)

func main() {
	// Mock login data
	loginData := broker.LoginData{
		Email:    "yoneclick@gmail.com",
		Password: "Junio020499",
	}

	broker.SetVerbose(true)

	// Mock host name
	hostName := "trade.exnova.com"

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

	client.OnTradeOpened(func(tradeData broker.TradeData) {
		fmt.Println(tradeData.TimeFrameInMinutes, " minute trade opened: ", tradeData.TradeID)
	})

	lastBalance := 0.0

	for {
		balance, _ := client.GetBalance(broker.BalanceTypeDemo)

		if balance.Amount == lastBalance {
			continue
		}

		fmt.Printf("Balance changed: %f\n", balance.Amount)

		lastBalance = balance.Amount
	}
}
