package broker_test

import (
	"testing"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/broker"
)

func TestBrokerClient(t *testing.T) {
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
		t.Fatalf("Login failed: %v\n", err)
	}

	err = client.ConnectSocket()
	if err != nil {
		t.Fatalf("Socket connection failed: %v\n", err)
	}

	for !client.IsReady() {
		time.Sleep(time.Second * 1)
	}

	t.Log("Client is ready")

	balances := client.GetBalances()
	if len(balances) == 0 {
		t.Fatal("Failed to retrieve balances")
	} else {
		for _, balance := range balances {
			t.Logf("%s: %f\n", balance.Currency, balance.Amount)
		}
	}

	client.OnTradeOpened(func(tradeData broker.TradeData) {
		t.Logf("%d minute trade opened: %d\n", tradeData.TimeFrameInMinutes, tradeData.TradeID)
	})

	newAssets, err := client.GetAssets(broker.AssetTypeDigital)
	if err != nil {
		t.Fatalf("Failed to get assets: %v\n", err)
	} else {
		assets := newAssets.WithoutNonTradable()

		for i, asset := range assets {
			t.Logf("%d: %d\n", i, asset.ActiveID)
		}
	}

	tradeType := broker.AssetTypeDigital
	tradeID := 0
	timeFrame := 1 // Time frame in minutes
	time.Sleep(time.Second * 1)

openTrade:
	// Replace parameters with appropriate values
	newTradeID, err := client.OpenTrade(tradeType, 100.0, broker.TradeDirectionCall, 1, timeFrame, broker.BalanceTypeDemo)
	if err != nil {
		t.Fatalf("Failed to open trade: %v\n", err)
	}

	// Ensure trade ID is valid
	if newTradeID == 0 {
		t.Fatal("Invalid trade ID")
	}

	tradeID = newTradeID

	_, err = client.CheckTradeResult(tradeID, timeFrame)
	if err != nil {
		t.Fatalf("Failed to check trade result: %v\n", err)
	}

	if tradeType == broker.AssetTypeDigital {
		tradeType = broker.AssetTypeBinary
		goto openTrade
	}
}
