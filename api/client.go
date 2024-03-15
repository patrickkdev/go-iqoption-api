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
	LoginData     *httpapi.LoginData
	Session       *httpapi.Session
	HostData      *data.Host
	WebSocket     *wsapi.Socket
	EventHandlers map[string]wsapi.EventCallback
	TimeSync      *wsapi.Timesync
	Balances      *wsapi.Balances
}

func NewBrokerClient(hostName string) *BrokerClient {
	return &BrokerClient{
		HostData:      data.GetHostData(hostName),
		EventHandlers: make(map[string]wsapi.EventCallback),
		TimeSync:      wsapi.NewTimesync(),
		Session:       httpapi.NewSession(),
	}
}

func (bC *BrokerClient) Login(email string, password string, token *string) (*BrokerClient, error) {
	data := httpapi.LoginData{
		Identifier: email,
		Password:   password,
		Token:      token,
	}

	err := httpapi.Login(bC.HostData.LoginURL, bC.Session, &data)

	if err != nil {
		return nil, err
	}

	bC.LoginData = &data

	return bC, nil
}

func (bC *BrokerClient) Logout() error {
	return httpapi.Logout(bC.HostData.LogoutURL, bC.Session)
}

func (bC *BrokerClient) ConnectSocket() error {
	socketConnection, err := wsapi.NewSocketConnection(bC.HostData.WSAPIURL)
	if err != nil {
		return err
	}

	bC.WebSocket = socketConnection

	// Handle heartbeat
	bC.WebSocket.AddEventListener("heartbeat", func(event []byte) {
		var heartbeatFromServer wsapi.Heartbeat
		json.Unmarshal(event, &heartbeatFromServer)

		wsapi.AnswerHeartbeat(bC.WebSocket, heartbeatFromServer, bC.TimeSync.GetServerTimestamp())
	})

	// Handle heartbeat
	bC.WebSocket.AddEventListener("timeSync", func(event []byte) {
		var timesyncEvent wsapi.TimesyncEvent
		json.Unmarshal(event, &timesyncEvent)

		bC.TimeSync.SetServerTimestamp(timesyncEvent.Msg)
	})

	// Handle authentication
	resp, err := wsapi.Authenticate(
		bC.WebSocket,
		bC.Session.SSID,
		int(bC.TimeSync.GetServerTimestamp()*1000),
		time.Now().Add(5*time.Second),
	)

	if err != nil {
		debug.IfVerbose.Println("Authentication error: ", err.Error())
		bC.WebSocket.Close()
		bC.WebSocket.WaitGroup.Done()
		return err
	}

	if !resp.Msg {
		error_ := fmt.Errorf("authentication error")
		debug.IfVerbose.Println(error_)
		debug.IfVerbose.PrintAsJSON(resp)
		return error_
	}

	fmt.Println("Authenticated successfully")

	return nil
}

func (bC *BrokerClient) GetCoreProfile() (*wsapi.CoreProfile, error) {
	return wsapi.GetCoreProfile(
		bC.WebSocket,
		bC.TimeSync.GetServerTimestamp(),
		time.Now().Add(5*time.Second),
	)
}

func (bC *BrokerClient) GetProfileClient(userId int) (*wsapi.UserProfileClient, error) {
	return wsapi.GetUserProfileClient(
		bC.WebSocket,
		userId,
		bC.TimeSync.GetServerTimestamp(),
		time.Now().Add(5*time.Second),
	)
}

func (bC *BrokerClient) GetBalances(shouldUpdate bool) (*wsapi.Balances, error) {
	if bC.Balances == nil || shouldUpdate {
		balances, err := wsapi.GetBalances(bC.WebSocket, bC.TimeSync.GetServerTimestamp())
		if err != nil {
			return nil, err
		}

		for _, balance := range *balances {
			bC.WebSocket.EmitEvent(wsapi.GetSubscriptionToPositionChanged(balance.UserID, balance.ID))
		}

		bC.Balances = balances
		return balances, nil
	}

	return bC.Balances, nil
}

func (bC *BrokerClient) OpenTrade(type_ wsapi.TradeType, amount float64, direction wsapi.TradeDirection, activeID int, duration int, balance wsapi.TradeBalance) (int, error) {
	balances, err := bC.GetBalances(false)
	if err != nil {
		return 0, err
	}

	targetBalance, err := balances.FindByType(balance)
	if err != nil {
		return 0, err
	}

	switch type_ {
	case wsapi.TradeTypeBinary:
		return wsapi.TradeBinary(
			bC.WebSocket,
			amount,
			direction,
			activeID,
			duration,
			targetBalance.ID,
			bC.TimeSync.GetServerTimestamp(),
		)
	case wsapi.TradeTypeDigital:
		return wsapi.TradeDigital(
			bC.WebSocket,
			amount,
			direction,
			activeID,
			duration,
			targetBalance.ID,
			bC.TimeSync.GetServerTimestamp(),
		)
	}

	return 0, nil
}
