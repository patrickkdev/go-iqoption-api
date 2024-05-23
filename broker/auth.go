package broker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/internal/brokerws"
	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
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
func (c *Client) ConnectSocket(connRetryCount ...int) error {
	maxRetryCount := 10
	retryCount := 0

	if len(connRetryCount) > 0 {
		retryCount = connRetryCount[0]
		fmt.Printf("Reconnecting user %s... (%d/%d)\n", c.session.LoginData.Email, retryCount, 10)
	}

	if retryCount >= maxRetryCount {
		return fmt.Errorf("max retry count reached for user %s", c.session.LoginData.Email)
	}

	// Try to persist connection
	reconnect := func() {
		// REFACTOR NEEDED
		
		return // MAKE NEXT CODE UNREACHABLE IN ORDER TO TEST FUNCIONALITY WITHOUT RECONNECT 
		
		// WE SUSPECT THAT THIS RECONNECT FUNCTION IS WHAT IS CAUSING THE SERVER TO SHUT DOWN
		time.Sleep(time.Minute)

		err := c.ConnectSocket(retryCount + 1)
		if err != nil {
			debug.IfVerbose.Println("Reconnect error", err.Error())
		}
	}

	newSocketConn, err := brokerws.NewSocketConnection(c.brokerEndpoints.WSAPIURL, c.defaultTimeoutDuration, reconnect)
	if err != nil {
		return err
	}

	c.ws = newSocketConn

	// REFACTOR NEEDED
	// It seems like answering heartbeats is not needed, so we remove it
	c.startAnsweringHeartBeats()

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
