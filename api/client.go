package api

import (
	"encoding/json"
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
		time.Now().Add(5 * time.Second),
	)

	debug.PrintAsJSON(resp)

	if err != nil {
		debug.IfVerbose.Println("Authentication error: ", err.Error())
		bC.WebSocket.Close()
		bC.WebSocket.WaitGroup.Done()
		return err
	}

	return nil
}

func (bC *BrokerClient) GetCoreProfile() (*wsapi.CoreProfile, error) {
	return wsapi.GetCoreProfile(
		bC.WebSocket, 
		bC.TimeSync.GetServerTimestamp(), 
		time.Now().Add(5 * time.Second),
	)
}

func (bC *BrokerClient) GetProfileClient(userId int) (*wsapi.UserProfileClient, error) {
	return wsapi.GetUserProfileClient(
		bC.WebSocket,
		userId,
		bC.TimeSync.GetServerTimestamp(),
		time.Now().Add(5 * time.Second),
	)
}