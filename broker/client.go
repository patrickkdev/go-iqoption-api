package broker

import (
	"fmt"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/internal/brokerws"
	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
	"github.com/patrickkdev/Go-IQOption-API/internal/types"
)

type Client struct {
	balances Balances

	session        Session
	brokerHostData types.Host

	ws *brokerws.WebSocket

	serverTimestamp int64

	// Default timeout duration for requests
	defaultTimeoutDuration time.Duration

	// Callbacks
	openTradeCallback     func(tradeData TradeData)
	onTradeClosedCallback map[int]func(tradeData TradeData)
}

// Create new client
func NewClient(loginData LoginData, hostName string, defaultTimeoutDuration time.Duration, ssid ...string) *Client {
	newClient := &Client{
		brokerHostData:         types.GetHostData(hostName),
		session:                NewSession(loginData),
		defaultTimeoutDuration: defaultTimeoutDuration,
		onTradeClosedCallback:  make(map[int]func(tradeData TradeData)),
	}

	if len(ssid) > 0 && ssid[0] != "" {
		newClient.session.SSID = ssid[0]
	}

	return newClient
}

// Logs in, gets the SSID and sets up the session
func (c *Client) Login() error {
	err := httpLogin(c.brokerHostData.LoginURL, &c.session)
	if err != nil {
		return err
	}

	return nil
}

// Logs out making the SSID invalid
func (c *Client) Logout() error {
	return httpLogout(c.brokerHostData.LogoutURL, &c.session)
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

	newSocketConn, err := brokerws.NewSocketConnection(c.brokerHostData.WSAPIURL, reconnect)
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

// Get all balances. These balances are kept updated by subscription to 'balance-changed' event
// Call GetUpdatedBalances if you need to really ensure that the balances are up-to-date
func (c *Client) GetBalances() Balances {
	return c.balances
}

// Get balance by type
func (c *Client) GetBalance(type_ BalanceType) (Balance, error) {
	return c.balances.FindByType(type_)
}

// Get session data
func (c *Client) GetSession() Session {
	return c.session
}

// Get broker host domain like 'iqoption.com'
func (c *Client) GetBrokerHost() string {
	return c.brokerHostData.Host
}

// Get timestamp that is kept updated by the 'timesync' event
func (c *Client) GetTimestamp() int64 {
	return c.serverTimestamp
}

// Get default timeout duration
func (c *Client) GetDefaultTimeoutDuration() time.Duration {
	return c.defaultTimeoutDuration
}

// Returns true if 'timeSync' event has been received and handled at least once
func (c *Client) IsReady() bool {
	return c.serverTimestamp != 0
}

// Returns true if the websocket connection not closed
func (c *Client) IsConnectionOK() bool {
	return !c.ws.Closed
}

// It is safe to call this function even if the connection is already closed
func (c *Client) CloseConnection() {
	if c.ws == nil || c.ws.Closed {
		return
	}

	c.ws.Close()
}

func (c *Client) handleTradeOpened(tradeData TradeData) {
	debug.IfVerbose.Printf("Trade %d opened", tradeData.TradeID)

	if c.openTradeCallback != nil {
		c.openTradeCallback(tradeData)
		debug.IfVerbose.Printf(" calling open trade callback\n")
	} else {
		debug.IfVerbose.Printf(" but there is no open trade callback set\n")
	}
}

func (c *Client) handleTradeClosed(tradeData TradeData) {
	debug.IfVerbose.Printf("Trade %d closed", tradeData.TradeID)

	callback, ok := c.onTradeClosedCallback[tradeData.TradeID]
	if !ok || callback == nil {
		debug.IfVerbose.Printf(" but there is no trade closed callback set for trade it\n")

		for index := range c.onTradeClosedCallback {
			fmt.Printf("Callback: %d\n", index)
		}

		return
	}

	debug.IfVerbose.Printf(" calling trade closed callback\n")
	callback(tradeData)
	delete(c.onTradeClosedCallback, tradeData.TradeID)
}

func (c *Client) getTimeout() time.Time {
	return time.Now().Add(c.defaultTimeoutDuration)
}
