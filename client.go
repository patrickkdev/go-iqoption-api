package broker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/brokerhttp"
	"github.com/patrickkdev/Go-IQOption-API/brokerws"

	"github.com/patrickkdev/Go-IQOption-API/internal/data"
	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
)

type Client struct {
	WS        *brokerws.Socket
	LoginData *brokerhttp.LoginData
	Balances  *brokerws.Balances

	hostData *data.Host
	session  *brokerhttp.Session

	timeSync *brokerws.Timesync

	// Timeout duration for requests
	timeoutDFR time.Duration

	openTradeCallback func(tradeData brokerws.TradeData)
}

func NewClient(loginData *brokerhttp.LoginData, hostName string, timeoutForRequests time.Duration, ssid ...string) *Client {
	newClient := &Client{
		LoginData:     loginData,
		hostData:      data.GetHostData(hostName),
		session:       brokerhttp.NewSession(),
		timeSync:      brokerws.NewTimesync(),
		timeoutDFR:    timeoutForRequests,
	}

	if len(ssid) > 0 && ssid[0] != "" {
		newClient.session.SSID = ssid[0]
	}

	return newClient
}

func (bC *Client) Login() (*Client, error) {
	err := brokerhttp.Login(bC.hostData.LoginURL, bC.session, bC.LoginData)

	if err != nil {
		return nil, err
	}

	return bC, nil
}

func (bC *Client) Logout() error {
	return brokerhttp.Logout(bC.hostData.LogoutURL, bC.session)
}

func (bC *Client) GetSSID() string {
	return bC.session.SSID
}

func (bC *Client) GetBrokerHost() string {
	return bC.hostData.Host
}

func (bC *Client) ConnectSocket() error {
	reconnect := func() {
		for {
			debug.IfVerbose.Println("Reconnecting...")
			err := bC.ConnectSocket()
			if err != nil {
				debug.IfVerbose.Println("Reconnect error: ", err.Error())
				time.Sleep(time.Second)
				continue
			}

			break
		}

		debug.IfVerbose.Println("Reconnected")
	}

	socketConnection, err := brokerws.NewSocketConnection(bC.hostData.WSAPIURL, reconnect)
	defer func () {
		debug.IfVerbose.Printf("ConnectSocket defer called with error: %v\n", err)

		if err != nil {
			if socketConnection != nil && !socketConnection.Closed {
				socketConnection.Close()
			}
			
			socketConnection = nil
		}
	}()

	if err != nil {
		return err
	}
	
	// Handle authentication
	resp, err := brokerws.Authenticate(
		socketConnection,
		bC.session.SSID,
		bC.getTimeout(),
	)

	if err != nil {
		return err
	}
	
	if !resp.Msg {
		err = fmt.Errorf("authentication error")
		debug.IfVerbose.PrintAsJSON(resp)
		return err
	}

	if socketConnection == nil {
		return fmt.Errorf("socket connection is nil")
	}
	
	bC.WS = socketConnection
	debug.IfVerbose.Println("Authenticated successfully")

	// region Handle Events

	// Heartbeat
	bC.WS.AddEventListener("heartbeat", func(event []byte) {
		var heartbeatFromServer brokerws.Heartbeat
		json.Unmarshal(event, &heartbeatFromServer)

		brokerws.AnswerHeartbeat(bC.WS, heartbeatFromServer, bC.timeSync.GetServerTimestamp())
	})

	// Timesync
	bC.WS.AddEventListener("timeSync", func(event []byte) {
		var timesyncEvent brokerws.TimesyncEvent
		json.Unmarshal(event, &timesyncEvent)

		bC.timeSync.SetServerTimestamp(timesyncEvent.Msg)
	})

	// Open Trade
	if bC.openTradeCallback != nil {
		bC.RemoveOpenTradeListener()
		bC.SetOpenTradeListener(bC.openTradeCallback)
	}

	// Get updated balances
	bC.GetBalances(true)
	return nil
}

func (bC *Client) SetOpenTradeListener(onOpenTrade func(tradeData brokerws.TradeData)) {
	bC.openTradeCallback = onOpenTrade
	brokerws.SetOpenTradeListener(bC.WS, onOpenTrade)
}

func (bC *Client) RemoveOpenTradeListener() {
	bC.openTradeCallback = nil
	brokerws.RemoveOpenTradeListener(bC.WS)
}

func (bC *Client) GetCoreProfile() (*brokerws.CoreProfile, error) {
	return brokerws.GetCoreProfile(
		bC.WS,
		bC.getTimeout(),
	)
}

func (bC *Client) GetProfileClient(userId int) (*brokerws.UserProfileClient, error) {
	return brokerws.GetUserProfileClient(
		bC.WS,
		userId,
		bC.getTimeout(),
	)
}

func (bC *Client) GetBalances(shouldUpdate bool) (*brokerws.Balances, error) {
	if bC.Balances == nil || shouldUpdate {
		balances, err := brokerws.GetBalances(bC.WS, bC.getTimeout())
		if err != nil {
			return nil, err
		}

		for _, balance := range *balances {
			subscriptionEvents := brokerws.GetSubscriptionsToPositionChangedEvent(balance.UserID, balance.ID)

			for _, event := range subscriptionEvents {
				debug.IfVerbose.Println("Subscribing to position changed: ")
				debug.IfVerbose.PrintAsJSON(event)
				brokerws.Emit(bC.WS, event)
			}
		}

		bC.Balances = balances
		return balances, nil
	}

	return bC.Balances, nil
}

func (bC *Client) GetCandles(count int, timeFrameInMinutes int, endtime int64, activeID int) (candles brokerws.Candles, err error) {
	return brokerws.GetCandles(bC.WS, count, timeFrameInMinutes, endtime, activeID, bC.getTimeout())
}

func (bC *Client) GetTopAssets(type_ brokerws.AssetType) (*brokerws.Assets, error) {
	return brokerws.GetTopAssets(bC.WS, type_, bC.getTimeout())
}

func (bC *Client) OpenTrade(type_ brokerws.AssetType, amount float64, direction brokerws.TradeDirection, activeID int, timeFrameInMinutes int, balance brokerws.BalanceType, shouldWaitForResult brokerws.TradeShouldWaitForResult) (int, bool, error) {
	balances, err := bC.GetBalances(false)
	if err != nil {
		return 0, false, err
	}

	targetBalance, err := balances.FindByType(balance)
	if err != nil {
		return 0, false, err
	}

	tradeID := 0
	win := false

	switch type_ {
	case brokerws.AssetTypeBinary:
		tradeID, err = brokerws.TradeBinary(
			bC.WS,
			amount,
			direction,
			activeID,
			timeFrameInMinutes,
			targetBalance.ID,
			bC.timeSync.GetServerTimestamp(),
			bC.getTimeout(),
		)
	case brokerws.AssetTypeDigital:
		tradeID, err = brokerws.TradeDigital(
			bC.WS,
			amount,
			direction,
			activeID,
			timeFrameInMinutes,
			targetBalance.ID,
			bC.timeSync.GetServerTimestamp(),
			bC.getTimeout(),
		)
	}

	if err != nil {
		return 0, false, err
	}

	if !bool(shouldWaitForResult) {
		return tradeID, false, nil
	}

	win, err = bC.CheckTradeWin(tradeID, type_, timeFrameInMinutes)
	return tradeID, win, err
}

func (bC *Client) CheckDigitalTradeResult(id int, timeFrameInMinutes int) (*brokerws.DigitalResult, bool, error) {
	return brokerws.CheckResultDigital(bC.WS, id, bC.getTimeout().Add(time.Minute*time.Duration(timeFrameInMinutes)))
}

func (bC *Client) CheckBinaryTradeResult(id int, timeFrameInMinutes int) (*brokerws.BinaryResult, bool, error) {
	return brokerws.CheckResultBinary(bC.WS, id, bC.getTimeout().Add(time.Minute*time.Duration(timeFrameInMinutes)))
}

func (bC *Client) CheckTradeWin(id int, type_ brokerws.AssetType, timeFrameInMinutes int) (bool, error) {
	win := false
	var err error = nil

	switch type_ {
	case brokerws.AssetTypeBinary:
		_, win, err = bC.CheckBinaryTradeResult(id, timeFrameInMinutes)
		return win, err
	case brokerws.AssetTypeDigital:
		_, win, err = bC.CheckDigitalTradeResult(id, timeFrameInMinutes)
		return win, err
	}

	return false, nil
}

func (bC *Client) getTimeout() time.Time {
	return time.Now().Add(bC.timeoutDFR)
}
