package api

import (
	"encoding/json"
	"fmt"
	"patrickkdev/Go-IQOption-API/data"
	"patrickkdev/Go-IQOption-API/debug"
	"patrickkdev/Go-IQOption-API/httpapi"
	"patrickkdev/Go-IQOption-API/wsapi"
	"time"
)

type BrokerClient struct {
	WS     							*wsapi.Socket
	LoginData     			*httpapi.LoginData
	Balances      			*wsapi.Balances
	
	hostData      			*data.Host
	session       			*httpapi.Session
	
	timeSync      			*wsapi.Timesync
	
	// Timeout duration for requests
	timeoutDFR 					time.Duration
	
	eventHandlers 			map[string]wsapi.EventCallback
}

func NewBrokerClient(hostName string, timeoutForRequests time.Duration) *BrokerClient {
	return &BrokerClient{
		hostData:      	data.GetHostData(hostName),
		session:       	httpapi.NewSession(),
		timeSync:      	wsapi.NewTimesync(),
		eventHandlers: 	make(map[string]wsapi.EventCallback),
		timeoutDFR: 		timeoutForRequests,
	}
}

func (bC *BrokerClient) Login(email string, password string, token *string) (*BrokerClient, error) {
	data := &httpapi.LoginData{
		Identifier: email,
		Password:   password,
		Token:      token,
	}

	err := httpapi.Login(bC.hostData.LoginURL, bC.session, data)

	if err != nil {
		return nil, err
	}

	bC.LoginData = data

	return bC, nil
}

func (bC *BrokerClient) Logout() error {
	return httpapi.Logout(bC.hostData.LogoutURL, bC.session)
}

func (bC *BrokerClient) ConnectSocket() error {
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

	socketConnection, err := wsapi.NewSocketConnection(bC.hostData.WSAPIURL, reconnect)
	if err != nil {
		return err
	}

	bC.WS = socketConnection

	// Handle authentication
	resp, err := wsapi.Authenticate(
		bC.WS,
		bC.session.SSID,
		bC.getTimeout(),
	)

	if err != nil {
		debug.IfVerbose.Println("Authentication error: ", err.Error())
		bC.WS.Close()
		bC.WS.WaitGroup.Done()
		return err
	}

	if !resp.Msg {
		error_ := fmt.Errorf("authentication error")
		debug.IfVerbose.Println(error_)
		debug.IfVerbose.PrintAsJSON(resp)
		return error_
	}

	fmt.Println("Authenticated successfully")

	// Handle heartbeat
	bC.WS.AddEventListener("heartbeat", func(event []byte) {
		var heartbeatFromServer wsapi.Heartbeat
		json.Unmarshal(event, &heartbeatFromServer)

		wsapi.AnswerHeartbeat(bC.WS, heartbeatFromServer, bC.timeSync.GetServerTimestamp())
	})

	// Handle timesync
	bC.WS.AddEventListener("timeSync", func(event []byte) {
		var timesyncEvent wsapi.TimesyncEvent
		json.Unmarshal(event, &timesyncEvent)

		bC.timeSync.SetServerTimestamp(timesyncEvent.Msg)
	})

	bC.GetBalances(true)

	return nil
}

func (bC *BrokerClient) GetCoreProfile() (*wsapi.CoreProfile, error) {
	return wsapi.GetCoreProfile(
		bC.WS,
		bC.getTimeout(),
	)
}

func (bC *BrokerClient) GetProfileClient(userId int) (*wsapi.UserProfileClient, error) {
	return wsapi.GetUserProfileClient(
		bC.WS,
		userId,
		bC.getTimeout(),
	)
}

func (bC *BrokerClient) GetBalances(shouldUpdate bool) (*wsapi.Balances, error) {
	if bC.Balances == nil || shouldUpdate {
		balances, err := wsapi.GetBalances(bC.WS, bC.getTimeout())
		if err != nil {
			return nil, err
		}

		for _, balance := range *balances {
			subscriptionEvents := wsapi.GetSubscriptionsToPositionChangedEvent(balance.UserID, balance.ID)

			for _, event := range subscriptionEvents {
				debug.IfVerbose.Println("Subscribing to position changed: ")
				debug.IfVerbose.PrintAsJSON(event)
				wsapi.Emit(bC.WS, event)
			}
		}

		bC.Balances = balances
		return balances, nil
	}

	return bC.Balances, nil
}

func (bC *BrokerClient) GetCandles(count int, timeFrameInMinutes int, endtime int64, activeID int) (candles wsapi.Candles, err error) {
	return wsapi.GetCandles(bC.WS, count, timeFrameInMinutes, endtime, activeID, bC.getTimeout())
}

func (bC *BrokerClient) OpenTrade(type_ wsapi.TradeType, amount float64, direction wsapi.TradeDirection, activeID int, timeFrameInMinutes int, balance wsapi.TradeBalance, waitForResult bool) (int, bool, error) {
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
	case wsapi.TradeTypeBinary:
		tradeID, err = wsapi.TradeBinary(
			bC.WS,
			amount,
			direction,
			activeID,
			timeFrameInMinutes,
			targetBalance.ID,
			bC.timeSync.GetServerTimestamp(),
			bC.getTimeout(),
		)
	case wsapi.TradeTypeDigital:
		tradeID, err = wsapi.TradeDigital(
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

	if !waitForResult {
		return tradeID, false, nil
	}
	
	win, err = bC.CheckTradeWin(tradeID, type_, timeFrameInMinutes)
	return tradeID, win, err
}

func (bC *BrokerClient) CheckDigitalTradeResult(id int, timeFrameInMinutes int) (*wsapi.DigitalResult, bool, error) {
	return wsapi.CheckResultDigital(bC.WS, id, bC.getTimeout().Add(time.Minute * time.Duration(timeFrameInMinutes)))
}

func (bC *BrokerClient) CheckBinaryTradeResult(id int, timeFrameInMinutes int) (*wsapi.BinaryResult, bool, error) {
	return wsapi.CheckResultBinary(bC.WS, id, bC.getTimeout().Add(time.Minute * time.Duration(timeFrameInMinutes)))
}

func (bC *BrokerClient) CheckTradeWin(id int, type_ wsapi.TradeType, timeFrameInMinutes int) (bool, error) {
	win := false
	var err error = nil

	switch type_ {
	case wsapi.TradeTypeBinary:
		_, win, err = bC.CheckBinaryTradeResult(id, timeFrameInMinutes)
		return win, err
	case wsapi.TradeTypeDigital:
		_, win, err = bC.CheckDigitalTradeResult(id, timeFrameInMinutes)
		return win, err
	}

	return false, nil
}

func (bC *BrokerClient) getTimeout() time.Time {
	return time.Now().Add(bC.timeoutDFR)
}
