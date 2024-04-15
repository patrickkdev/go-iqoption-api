package broker

import (
	"fmt"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/btypes"
	"github.com/patrickkdev/Go-IQOption-API/internal/brokerhttp"
	"github.com/patrickkdev/Go-IQOption-API/internal/brokerws"

	"github.com/patrickkdev/Go-IQOption-API/internal/data"
	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
)

type Client struct {
	balances btypes.Balances

	session  brokerhttp.Session
	hostData data.Host

	ws *brokerws.Socket

	time btypes.Timesync

	// Default timeout duration for requests
	defaultTimeoutDuration time.Duration

	// Callbacks
	openTradeCallback     func(tradeData btypes.TradeData)
	onTradeClosedCallback map[int]func(tradeData btypes.TradeData)
}

// Create new client
func NewClient(loginData brokerhttp.LoginData, hostName string, timeoutForRequests time.Duration, ssid ...string) *Client {
	newClient := &Client{
		hostData:               data.GetHostData(hostName),
		session:                brokerhttp.NewSession(loginData),
		time:                   brokerws.NewTimesync(),
		defaultTimeoutDuration: timeoutForRequests,
		onTradeClosedCallback:  make(map[int]func(tradeData btypes.TradeData)),
	}

	if len(ssid) > 0 && ssid[0] != "" {
		newClient.session.SSID = ssid[0]
	}

	return newClient
}

// Get all balances. These balances are kept updated by subscription to 'balance-changed' event
// Call GetUpdatedBalances if you need to really ensure that the balances are up-to-date
func (bC *Client) GetBalances() btypes.Balances {
	return bC.balances
}

// Get balance by type
func (bC *Client) GetBalance(type_ btypes.BalanceType) (btypes.Balance, error) {
	return bC.balances.FindByType(type_)
}

// Get session data
func (bC *Client) GetSession() brokerhttp.Session {
	return bC.session
}

// Get host data
func (bC *Client) GetHostData() data.Host {
	return bC.hostData
}

// Get client time data that is kept updated by subscription to 'timesync' event
func (bC *Client) GetTime() btypes.Timesync {
	return bC.time
}

// Get default timeout duration
func (bC *Client) GetDefaultTimeoutDuration() time.Duration {
	return bC.defaultTimeoutDuration
}

func (bC *Client) IsReady() bool {
	return bC.time.GetServerTimestamp() != 0
}

// Logs in, gets the SSID and sets up the session
func (bC *Client) Login() (*Client, error) {
	err := brokerhttp.Login(bC.hostData.LoginURL, &bC.session)
	if err != nil {
		return nil, err
	}

	return bC, nil
}

// Logs out making the SSID invalid
func (bC *Client) Logout() error {
	return brokerhttp.Logout(bC.hostData.LogoutURL, &bC.session)
}

// Connect to broker websocket and authenticates
// Is always trying to reconnect if connection is lost
// Automatically subscribes to default events like 'balance-changed', 'trade-changed', 'timesync' and 'position-changed'
func (bC *Client) ConnectSocket() error {
	// Persist connection
	reconnect := func() {
		for {
			debug.IfVerbose.Println("Reconnecting...")
			err := bC.ConnectSocket()
			if err != nil {
				debug.IfVerbose.Println("Reconnect error: ", err.Error())
				time.Sleep(time.Second)
				continue
			}

			break
		}

		debug.IfVerbose.Println("Reconnected")
	}

	socketConnection, err := brokerws.NewSocketConnection(bC.hostData.WSAPIURL, reconnect)
	defer func() {
		debug.IfVerbose.Printf("ConnectSocket defer called with error: %v\n", err)

		if err != nil {
			if socketConnection != nil && !socketConnection.Closed {
				socketConnection.Close()
			}

			socketConnection = nil
		}
	}()

	if err != nil {
		return err
	}

	// Handle authentication
	resp, err := socketConnection.Authenticate(
		bC.session.SSID,
		bC.getTimeout(),
	)

	if err != nil {
		return err
	}

	if !resp.Msg {
		err = fmt.Errorf("authentication error")
		debug.IfVerbose.PrintAsJSON(resp)
		return err
	}

	if socketConnection == nil {
		return fmt.Errorf("socket connection is nil")
	}

	// Successfully connected
	bC.ws = socketConnection

	// Get updated balances
	bC.balances, err = bC.GetUpdatedBalances()
	if err != nil {
		return err
	}

	// Heartbeat
	bC.ws.OnHeartBeat(func(heartbeat btypes.Heartbeat) {
		bC.ws.AnswerHeartBeat(heartbeat, bC.time.GetServerTimestamp())
	})

	// Timesync
	bC.ws.OnTimeSync(func(newTimestamp int64) {
		bC.time.SetServerTimestamp(newTimestamp)
	})

	// Trade
	bC.ws.SubscribeToTradeChanges(bC.balances)
	bC.ws.OnTradeChanged(brokerws.TradeChangedCallbacks{
		OnTradeOpenedCallback: bC.whenATradeOpens,
		OnTradeClosedCallback: bC.whenATradeCloses,
	})

	// Update stored balance when it changes
	bC.ws.OnBalanceChanged(func(update btypes.BalanceUpdate) {
		bC.balances.Update(update)
	})

	return nil
}

// Returns true if the websocket connection not closed
func (bC *Client) IsConnectionOK() bool {
	return !bC.ws.Closed
}

// It is safe to call this function even if the connection is already closed
func (bC *Client) CloseConnection() {
	if bC.ws == nil || bC.ws.Closed {
		return
	}

	bC.ws.Close()
}

// Get core profile data
// This is mostly used to get clientSessionID in order to call GetProfile
func (bC *Client) GetCoreProfile() (btypes.CoreProfile, error) {
	return bC.ws.GetCoreProfile(bC.getTimeout())
}

// Get user profile data
// You need to pass clientSessionID in order to call this function. Get it from GetCoreProfile function
func (bC *Client) GetProfile(clientSessionID int) (btypes.UserProfileClient, error) {
	return bC.ws.GetUserProfileClient(clientSessionID, bC.getTimeout())
}

// Although balances are kept updated by subscription to 'balance-changed' event,
// this function can be used to force an update and return the latest balances
func (bC *Client) GetUpdatedBalances() (btypes.Balances, error) {
	balances, err := bC.ws.GetBalances(bC.getTimeout())
	if err != nil {
		return nil, err
	}

	bC.balances = balances

	return balances, nil
}

// Gets OHLC candle data for a given timeframe and active id
func (bC *Client) GetCandles(count int, timeFrameInMinutes int, endtime int64, activeID int) (candles btypes.Candles, err error) {
	return bC.ws.GetCandles(count, timeFrameInMinutes, endtime, activeID, bC.getTimeout())
}

// Gets available assets (or 'pairs' like 'EUR/USD') for a given asset type like 'binary-option' or 'digital-option'
// Pairs returned are not garanteed to be tradable.
// Calling .FilterOpen() on the returned assets will TRY to filter out non-tradable pairs
func (bC *Client) GetAssets(type_ btypes.AssetType) (btypes.Assets, error) {
	return bC.ws.GetTopAssets(type_, bC.getTimeout())
}

// Immediately opens a new trade if params are valid and the asset is available
func (bC *Client) OpenTrade(type_ btypes.AssetType, amount float64, direction btypes.TradeDirection, activeID int, timeFrameInMinutes int, balance btypes.BalanceType) (int, error) {
	tradeID := 0

	targetBalance, err := bC.balances.FindByType(balance)
	if err != nil {
		return tradeID, err
	}

	switch type_ {
	case btypes.AssetTypeBinary:
		tradeID, err = bC.ws.TradeBinary(
			amount,
			direction,
			activeID,
			timeFrameInMinutes,
			targetBalance.ID,
			bC.time.GetServerTimestamp(),
			bC.getTimeout(),
		)
	case btypes.AssetTypeDigital:
		tradeID, err = bC.ws.TradeDigital(
			amount,
			direction,
			activeID,
			timeFrameInMinutes,
			targetBalance.ID,
			bC.time.GetServerTimestamp(),
			bC.getTimeout(),
		)
	}

	if err != nil {
		return tradeID, err
	}

	return tradeID, nil
}

// Waits until the trade is closed or timed out and returns the trade data
func (bC *Client) CheckTradeResult(tradeID int, timeFrameInMinutes int) (btypes.TradeData, error) {
	result := btypes.TradeData{}

	if tradeID == 0 {
		return result, fmt.Errorf("invalid trade id")
	}

	timedOut := time.Now().After(bC.getTimeout().Add(time.Duration(timeFrameInMinutes)))
	var err error = nil

	debug.IfVerbose.Println("Setting onTradeClosedCallback", tradeID)
	bC.onTradeClosed(tradeID, func(tradeData btypes.TradeData) {
		debug.IfVerbose.Println("CheckTradeResult :: Trade closed callback", tradeID)
		result = tradeData
	})

	debug.IfVerbose.Println("Waiting for trade result", tradeID)
	for result.TradeID != tradeID && !timedOut {
		time.Sleep(time.Second * 1)
	}

	if timedOut {
		err = fmt.Errorf("timed out waiting for trade result")
		delete(bC.onTradeClosedCallback, tradeID)
	}

	return result, err
}

// If called more than once, the previous callback will not be called
func (bC *Client) OnTradeOpened(callback func(tradeData btypes.TradeData)) {
	bC.openTradeCallback = callback
}

func (bC *Client) onTradeClosed(tradeID int, callback func(tradeData btypes.TradeData)) {
	bC.onTradeClosedCallback[tradeID] = callback
}

func (bC *Client) whenATradeOpens(tradeData btypes.TradeData) {
	debug.IfVerbose.Printf("Trade %d opened", tradeData.TradeID)

	if bC.openTradeCallback != nil {
		bC.openTradeCallback(tradeData)
		debug.IfVerbose.Printf(" calling open trade callback\n")
	} else {
		debug.IfVerbose.Printf(" but there is no open trade callback set\n")
	}
}

func (bC *Client) whenATradeCloses(tradeData btypes.TradeData) {
	debug.IfVerbose.Printf("Trade %d closed\n", tradeData.TradeID)

	callback, ok := bC.onTradeClosedCallback[tradeData.TradeID]
	if !ok || callback == nil {
		debug.IfVerbose.Printf(" but there is no trade closed callback set\n")

		for index := range bC.onTradeClosedCallback {
			fmt.Printf("Callback: %d\n", index)
		}

		return
	}

	debug.IfVerbose.Printf(" calling trade closed callback\n")
	callback(tradeData)
	delete(bC.onTradeClosedCallback, tradeData.TradeID)
}

func (bC *Client) getTimeout() time.Time {
	return time.Now().Add(bC.defaultTimeoutDuration)
}
