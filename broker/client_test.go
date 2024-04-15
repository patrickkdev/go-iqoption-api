package broker_test

import (
	"testing"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/broker"
	"github.com/patrickkdev/Go-IQOption-API/btypes"
	"github.com/patrickkdev/Go-IQOption-API/internal/brokerhttp"
)

func TestBrokerClient(t *testing.T) {
	// Mock login data
	loginData := brokerhttp.LoginData{
		Email:    "patrickfxtrader8q@gmail.com",
		Password: "YOUTAP2019",
	}

	// Mock host name
	hostName := "iqoption.com"

	// Timeout for requests
	timeout := time.Second * 10
	client := broker.NewClient(loginData, hostName, timeout)

	_, err := client.Login()
	if err != nil {
		t.Fatalf("Login failed: %v\n", err)
	}

	err = client.ConnectSocket()
	if err != nil {
		t.Fatalf("Socket connection failed: %v\n", err)
	}

	for !client.IsConnectionOK() {
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

	client.OnTradeOpened(func(tradeData btypes.TradeData) {
		t.Logf("%d minute trade opened: %d\n", tradeData.TimeFrameInMinutes, tradeData.TradeID)
	})

	newAssets, err := client.GetAssets(btypes.AssetTypeDigital)
	if err != nil {
		t.Fatalf("Failed to get assets: %v\n", err)
	} else {
		assets := newAssets.FilterOpen()

		for i, asset := range assets {
			t.Logf("%d: %d\n", i, asset.ActiveID)
		}
	}

	tradeType := btypes.AssetTypeDigital
	tradeID := 0
	timeFrame := 1 // Time frame in minutes
	time.Sleep(time.Second * 1)

openTrade:
	// Replace parameters with appropriate values
	newTradeID, err := client.OpenTrade(tradeType, 100.0, btypes.TradeDirectionCall, 1, timeFrame, btypes.BalanceTypeDemo)
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

	if tradeType == btypes.AssetTypeDigital {
		tradeType = btypes.AssetTypeBinary
		timeFrame = 2
		goto openTrade
	}
}
