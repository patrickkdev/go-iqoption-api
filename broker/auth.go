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
func (c *Client) ConnectSocket() error {
	// Persist connection
	reconnect := func() {
		for {
			debug.IfVerbose.Println("Reconnecting...")
			err := c.ConnectSocket()
			if err != nil {
				debug.IfVerbose.Println("Reconnect error: ", err.Error())
				time.Sleep(time.Second)
				continue
			}

			break
		}

		debug.IfVerbose.Println("Reconnected")
	}

	newSocketConn, err := brokerws.NewSocketConnection(c.brokerEndpoints.WSAPIURL, reconnect)
	if err != nil {
		return err
	}

	c.ws = newSocketConn

	c.startAnsweringHeartBeats()
	c.keepServerTimestampUpdated()

	// Handle authentication
	resp, err := c.authenticate()
	if err != nil {
		return err
	}

	if !resp.Msg {
		return fmt.Errorf("authentication error")
	}

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