package broker

import (
	"fmt"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
)

type TradeDirection string

const (
	TradeDirectionCall TradeDirection = "call"
	TradeDirectionPut  TradeDirection = "put"
)

type TradeData struct {
	Status             string
	TradeID            int
	Type               AssetType
	Direction          TradeDirection
	TimeFrameInMinutes int
	ActiveID           int
	Amount             float64
	Win                bool
	OpenTime           int
	Profit             float64
}

// Waits until the trade is closed or timed out and returns the trade data
func (c *Client) CheckTradeResult(tradeID int, timeFrameInMinutes int) (TradeData, error) {
	result := TradeData{}

	if tradeID == 0 {
		return result, fmt.Errorf("invalid trade id")
	}

	timedOut := time.Now().After(c.getTimeout().Add(time.Duration(timeFrameInMinutes)))
	var err error = nil

	c.onTradeClosed(tradeID, func(tradeData TradeData) {
		result = tradeData
	})

	debug.IfVerbose.Println("Waiting for trade result", tradeID)
	for result.TradeID != tradeID && !timedOut {
		time.Sleep(time.Second * 1)
	}

	if timedOut {
		err = fmt.Errorf("timed out waiting for trade result")
		delete(c.onTradeClosedCallback, tradeID)
	}

	return result, err
}
