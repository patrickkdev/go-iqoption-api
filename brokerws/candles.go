package brokerws

import (
	"fmt"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

type getCandlesResponse struct {
	Msg struct {
		Candles []Candle `json:"candles"`
	} `json:"msg"`
	Name      string `json:"name"`
	RequestID string `json:"request_id"`
	Status    int    `json:"status"`
}

type Candle struct {
	At     int64   `json:"at"`
	Close  float64 `json:"close"`
	From   int     `json:"from"`
	ID     int     `json:"id"`
	Max    float64 `json:"max"`
	Min    float64 `json:"min"`
	Open   float64 `json:"open"`
	To     int     `json:"to"`
	Volume int     `json:"volume"`
}

type Candles []Candle

func (candles Candles) GetLast() (candle Candle) {
	return candles[len(candles)-1]
}

func GetCandles(ws *Socket, count int, timeFrameInMinutes int, endtime int64, activeID int, timeout time.Time) (candles Candles, err error) {
	msg := map[string]interface{}{
		"name":    "get-candles",
		"version": "2.0",
		"body": map[string]interface{}{
			"active_id":           activeID,
			"split_normalization": true,
			"size":                timeFrameInMinutes * 60,
			"to":                  endtime,
			"count":               count,
		},
	}

	requestEvent := &RequestEvent{
		Name: "sendMessage",
		Msg:  msg,
	}

	resp, err := EmitWithResponse(ws, requestEvent, "candles", timeout)
	if err != nil {
		return nil, err
	}

	responseEvent, err := tjson.Unmarshal[getCandlesResponse](resp)
	if err != nil {
		return nil, err
	}

	debug.IfVerbose.PrintAsJSON(responseEvent)

	if responseEvent.Status != 2000 {
		return nil, fmt.Errorf("error getting candles")
	}

	return responseEvent.Msg.Candles, nil
}
