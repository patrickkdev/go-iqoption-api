package broker

import (
	"encoding/json"
	"fmt"

	"github.com/patrickkdev/Go-IQOption-API/internal/brokerws"
)

type LoginData struct {
	Email    string  `json:"identifier"`
	Password string  `json:"password"`
	Token    *string `json:"token,omitempty"`
}

// Logs in, gets the SSID and sets up the session
// Should be called after NewClient function
func (c *Client) Login() error {
	resp, err := c.session.PostFromStruct(c.brokerEndpoints.LoginURL, c.session.LoginData, nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&c.session)

	if err != nil {
		return err
	}

	return nil
}

// Logs out making the SSID invalid
func Logout(url string, session *Session) error {
	resp, err := session.PostFromStruct(url, nil, nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

// Connect to broker websocket and authenticates
// Is always trying to reconnect if connection is lost
// Automatically subscribes to default events like 'balance-changed', 'trade-changed', 'timesync' and 'position-changed'
func (c *Client) ConnectSocket() error {
	maxRetryCount := 10
	retryCount := 0

	if retryCount >= maxRetryCount {
		return fmt.Errorf("max retry count reached for user %s", c.session.LoginData.Email)
	}

	newSocketConn, err := brokerws.NewSocketConnection(c.brokerEndpoints.WSAPIURL, c.defaultTimeoutDuration)
	if err != nil {
		return err
	}

	c.ws = newSocketConn

	c.keepServerTimestampUpdated()

	// Handle authentication
	_, err = c.authenticate()
	if err != nil {
		return err
	}

	// Reset retry count if authentication is successful
	retryCount = 0

	// Get updated balances
	_, err = c.GetUpdatedBalances()
	if err != nil {
		return err
	}

	c.keepBalancesUpdated()

	c.subscribeToTradeChanges()
	c.onTradeChanged(TradeChangedCallbacks{
		onTradeOpened: c.handleTradeOpened,
		onTradeClosed: c.handleTradeClosed,
	})

	return nil
}

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

	if !responseEvent.Msg {
		return nil, fmt.Errorf("authentication failed")
	}

	return &responseEvent, nil
}

