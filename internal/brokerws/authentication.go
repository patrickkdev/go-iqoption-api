package brokerws

import (
	"time"

	"github.com/patrickkdev/Go-IQOption-API/btypes"
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

func (ws *Socket) Authenticate(ssid string, timeout time.Time) (*btypes.AuthenticationResponse, error) {
	requestEvent := &btypes.RequestEvent{
		Name: "authenticate",
		Msg: map[string]interface{}{
			"ssid":              ssid,
			"protocol":          3,
			"client_session_id": "",
			"session_id":        "",
		},
	}

	resp, err := ws.EmitWithResponse(requestEvent, "authenticated", timeout)
	if err != nil {
		return nil, err
	}

	responseEvent, err := tjson.Unmarshal[btypes.AuthenticationResponse](resp)
	if err != nil {
		return nil, err
	}

	return &responseEvent, nil
}
