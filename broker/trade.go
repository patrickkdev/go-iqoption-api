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

// CheckTradeResult waits until the trade is closed or timed out and returns the trade data
func (c *Client) CheckTradeResult(tradeID int, timeFrameInMinutes int) (TradeData, error) {
	result := TradeData{}

	if tradeID == 0 {
		return result, fmt.Errorf("invalid trade id")
	}

	timeout := c.getTimeout() + time.Minute * time.Duration(timeFrameInMinutes) + time.Minute

	ctx, cancel := context.WithDeadline(context.Background(), timeout)
	defer cancel()

	resultChan := make(chan TradeData)
	errorChan := make(chan error)

	c.onTradeClosed(tradeID, func(tradeData TradeData) {
		select {
		case resultChan <- tradeData:
		case <-ctx.Done():
		}
	})

	go func() {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			errorChan <- fmt.Errorf("timed out waiting for trade result")
		}
	}()

	debug.IfVerbose.Println("Waiting for trade result", tradeID)

	select {
	case result = <-resultChan:
		debug.IfVerbose.Println("Received trade result for", tradeID)
	case err := <-errorChan:
		delete(c.onTradeClosedCallback, tradeID)
		return result, err
	}

	delete(c.onTradeClosedCallback, tradeID)
	return result, nil
}
