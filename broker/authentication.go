package broker

import (
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
	"github.com/patrickkdev/Go-IQOption-API/internal/types"
)

func (c *Client) authenticate() (*types.AuthenticationResponse, error) {
	requestEvent := requestEvent{
		Name: "authenticate",
		Msg: map[string]interface{}{
			"ssid":              c.session.SSID,
			"protocol":          3,
			"client_session_id": "",
			"session_id":        "",
		},
	}.WithRandomRequestId()

	resp, err := c.ws.EmitWithResponse(requestEvent, "authenticated", c.getTimeout())
	if err != nil {
		return nil, err
	}

	responseEvent, err := tjson.Unmarshal[types.AuthenticationResponse](resp)
	if err != nil {
		return nil, err
	}

	return &responseEvent, nil
}
