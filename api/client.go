package api

import (
	"sync"

	"patrickkdev/Go-IQOption-API/data"
	"patrickkdev/Go-IQOption-API/httpapi"
	"patrickkdev/Go-IQOption-API/utils"
	"patrickkdev/Go-IQOption-API/wsapi"
)

type BrokerClient struct {
	LoginData *httpapi.LoginData
	Session *httpapi.Session
	HostData  *data.Host
	WebSocket *wsapi.WSocket
	EventHandlers map[string]wsapi.WSEventCallback
}

func NewBrokerClient(hostName string) *BrokerClient { 
	return &BrokerClient{
		HostData: data.GetHostData(hostName),
		Session: &httpapi.Session{
			Headers: nil,
			Cookie: "",
		},
		EventHandlers: make(map[string]wsapi.WSEventCallback),
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

func (bC *BrokerClient) ConnectSocket() (*sync.WaitGroup, error) {
	socketConnection, err := wsapi.NewSocketConnection("ws://localhost:8080") //(bC.HostData.WSAPIURL)
	if err != nil {
		return nil, err
	}

	bC.WebSocket = socketConnection

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go bC.WebSocket.Listen(wg, bC.HandleEvent)

	return wg, nil
}

func (bC *BrokerClient) Subscribe(name string, callback wsapi.WSEventCallback) {
	bC.EventHandlers[name] = callback
}

func (bC *BrokerClient) SendEvent(event interface{}) {
	bC.WebSocket.Write(event)
}

func (bC *BrokerClient) HandleEvent(event wsapi.WSEvent) {
	eventName, ok := event["name"].(string)
	if !ok {
		utils.PrintlnIfVerbose("no event name")
		return
	}

	callback, ok := bC.EventHandlers[eventName]
	if !ok {
		utils.PrintlnIfVerbose("no event callback")
		return
	}

	callback(event)
}

func (bC *BrokerClient) Close() {
	bC.WebSocket.Close()
}