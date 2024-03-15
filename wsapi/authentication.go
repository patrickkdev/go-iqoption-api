package wsapi

import (
	"time"

	"github.com/patrickkdev/Go-IQOption-API/tjson"
)

type AuthenticationResponse struct {
	Name            string `json:"name"`
	Msg             bool   `json:"msg"`
	ClientSessionID string `json:"client_session_id"`
	RequestID       string `json:"request_id"`
}

func Authenticate(ws *Socket, ssid string, timeout time.Time) (*AuthenticationResponse, error) {
	requestEvent := &RequestEvent{
		Name: "authenticate",
		Msg: map[string]interface{}{
			"ssid":              ssid,
			"protocol":          3,
			"client_session_id": "",
			"session_id":        "",
		},
	}

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
