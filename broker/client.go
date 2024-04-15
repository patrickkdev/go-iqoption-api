package broker

import (
	"time"

	"github.com/patrickkdev/Go-IQOption-API/internal/brokerws"
)

type Client struct {
	balances 								Balances

	session        					Session
	brokerEndpoints 				brokerEndpoints

	ws 											*brokerws.WebSocket

	serverTimestamp			 		int64

	// Default timeout duration for requests
	defaultTimeoutDuration 	time.Duration

	// Callbacks
	openTradeCallback     	func(tradeData TradeData)
	onTradeClosedCallback 	map[int]func(tradeData TradeData)
}

// Create new client
// You can optionally pass SSID to resume session with the broker ws. If omitted, a new ssid will be created when you call Login
// Recommended timeout duration is 10 seconds (time.Second * 10)
func NewClient(loginData LoginData, brokerDomain string, defaultTimeoutDuration time.Duration, ssid ...string) *Client {
	newClient := &Client{
		brokerEndpoints:         		getEndpointsByBrokerDomain(brokerDomain),
		session:                		NewSession(loginData),
		defaultTimeoutDuration: 		defaultTimeoutDuration,
		onTradeClosedCallback:  		make(map[int]func(tradeData TradeData)),
	}

	if len(ssid) > 0 && ssid[0] != "" {
		newClient.session.SSID = ssid[0]
	}

	return newClient
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
func (c *Client) GetBrokerDomain() string {
	return c.brokerEndpoints.Domain
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
	return c.ws.IsConnectionOK()
}

// It is safe to call this function even if the connection is already closed
func (c *Client) CloseConnection() {
	if c.ws == nil || c.ws.Closed {
		return
	}

	c.ws.Close()
}

func (c *Client) getTimeout() time.Time {
	return time.Now().Add(c.defaultTimeoutDuration)
}
