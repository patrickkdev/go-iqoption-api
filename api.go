package IQApi

import (
	"github.com/patrickkdev/Go-IQOption-API/data"
	httpapi "github.com/patrickkdev/Go-IQOption-API/http"
)

type BrokerClient struct {
	HostData  *data.Host
	LoginData httpapi.LoginData
}

func NewBrokerClient(hostName string) *BrokerClient {
	hostData := Data.GetHostData(hostName)
	return &BrokerClient{
		HostData: hostData,
	}
}

func (bC *BrokerClient) Login(email string, password string, token *string) (*BrokerClient, error) {
	data := httpapi.LoginData{
		Identifier: email,
		Password:   password,
		Token:      token,
	}

	err := httpapi.Login(bC.HostData.LoginURL, data)

	if err != nil {
		return nil, err
	}

	bC.LoginData = data

	return bC, nil
}
