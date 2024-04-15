package brokerws

import (
	"fmt"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/btypes"
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

type getCandlesResponse struct {
	Msg struct {
		Candles []btypes.Candle `json:"candles"`
	} `json:"msg"`
	Name      string `json:"name"`
	RequestID string `json:"request_id"`
	Status    int    `json:"status"`
}

func (ws *Socket) GetCandles(count int, timeFrameInMinutes int, endtime int64, activeID int, timeout time.Time) (candles btypes.Candles, err error) {
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

	requestEvent := &btypes.RequestEvent{
		Name: "sendMessage",
		Msg:  msg,
	}

	resp, err := ws.EmitWithResponse(requestEvent, "candles", timeout)
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
