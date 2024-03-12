package wsapi

import (
	"fmt"
	"patrickkdev/Go-IQOption-API/tjson"
	"time"
)

type AuthenticationResponse struct {
	Name            string `json:"name"`
	Msg             bool   `json:"msg"`
	ClientSessionID string `json:"client_session_id"`
	RequestID       string `json:"request_id"`
}

func Authenticate(ws *Socket, ssid string, serverTimeStamp int, timeout time.Time) (*AuthenticationResponse, error) {
	requestEvent := NewEvent(
		"authenticate",
		map[string]interface{}{
			"ssid":              ssid,
			"protocol":          3,
			"client_session_id": "",
			"session_id":        "",
		},
		fmt.Sprint(serverTimeStamp),
	)

	resp, err := EmitWithResponse(ws, requestEvent, "authenticated", timeout)
	if err != nil {
		return nil, err
	}

	responseEvent, err := tjson.Unmarshal[AuthenticationResponse](resp)
	if err != nil {
		return nil, err
	}

	return &responseEvent, nil
}