package broker

import (
	"fmt"

	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

type getCandlesResponse struct {
	Msg struct {
		Candles []Candle `json:"candles"`
	} `json:"msg"`
	Status    int    `json:"status"`
	// Name      string `json:"name"`
	// RequestID string `json:"request_id"`
}

type Candle struct {
	Open   float64 `json:"open"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Close  float64 `json:"close"`
	// At     int64   `json:"at"`
	// From   int     `json:"from"`
	// ID     int     `json:"id"`
	// To     int     `json:"to"`
	// Volume float64 `json:"volume"`
}

type Candles []Candle

// Gets OHLC candle data for a given timeframe and active id
func (c *Client) GetCandles(count int, timeFrameInMinutes int, endtime int64, activeID int) (candles Candles, err error) {
	requestEvent := requestEvent{
		Name: "sendMessage",
		Msg: map[string]interface{}{
			"name":    "get-candles",
			"version": "2.0",
			"body": map[string]interface{}{
				"active_id":           activeID,
				"split_normalization": true,
				"size":                timeFrameInMinutes * 60,
				"to":                  endtime,
				"count":               count,
			},
		},
	}.WithRandomRequestId()

	resp, err := c.ws.EmitWithResponse(requestEvent, "candles", c.getTimeout())
	if err != nil {
		return nil, err
	}

	responseEvent, err := tjson.Unmarshal[getCandlesResponse](resp)
	if err != nil {
		return nil, err
	}

	if responseEvent.Status != 2000 {
		return nil, fmt.Errorf("error getting candles")
	}

	return responseEvent.Msg.Candles, nil
}

func (candles Candles) GetLast() (candle Candle) {
	return candles[len(candles)-1]
}
