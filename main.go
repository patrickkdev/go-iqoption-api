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

	client.OnTradeOpened(func(tradeData broker.TradeData) {
		fmt.Println(tradeData.TimeFrameInMinutes, " minute trade opened: ", tradeData.TradeID)
	})

	for {
		balance, _ := client.GetBalance(broker.BalanceTypeDemo)

		fmt.Printf("Demo Balance: %f\n", balance.Amount)
		time.Sleep(time.Second * 1)
	}
}
