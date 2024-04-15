package broker

import (
	"strings"

	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
	"github.com/patrickkdev/Go-IQOption-API/internal/types"
)

type TradeChangedCallbacks struct {
	onTradeOpened func(tradeData TradeData)
	onTradeClosed func(tradeData TradeData)
}

// If called more than once, the previous callback will not be called
func (c *Client) OnTradeOpened(callback func(tradeData TradeData)) {
	c.openTradeCallback = callback
}

func (c *Client) onTradeClosed(tradeID int, callback func(tradeData TradeData)) {
	c.onTradeClosedCallback[tradeID] = callback
}

func (c *Client) onTradeChanged(callbacks TradeChangedCallbacks) {
	lastTradeID := make(map[AssetType]int)

	c.ws.AddEventListener("position-changed", func(event []byte) {
		eventString := string(event)
		tradeData := TradeData{ActiveID: 0}

		if strings.Contains(eventString, "binary_options_option_changed1") {
			res, err := tjson.Unmarshal[types.BinaryTradeData](event)
			if err != nil {
				return
			}

			tradeData.Status = res.Msg.Status
			tradeData.TradeID = res.Msg.ExternalID
			tradeData.Type = AssetTypeBinary
			tradeData.ActiveID = res.Msg.ActiveID
			tradeData.TimeFrameInMinutes = max((res.Msg.RawEvent.BinaryOptionsOptionChanged1.ExpirationTime-res.Msg.RawEvent.BinaryOptionsOptionChanged1.OpenTime)/60, 1)
			tradeData.Amount = res.Msg.RawEvent.BinaryOptionsOptionChanged1.Amount
			tradeData.Direction = TradeDirection(res.Msg.RawEvent.BinaryOptionsOptionChanged1.Direction)
			tradeData.Win = res.Msg.CloseReason == "win"
			tradeData.OpenTime = res.Msg.OpenTime
			tradeData.Profit = res.Msg.Pnl

		} else if strings.Contains(eventString, "digital_options_position_changed1") {
			res, err := tjson.Unmarshal[types.DigitalTradeData](event)
			if err != nil {
				return
			}

			tradeData.Status = res.Msg.Status
			tradeData.TradeID = res.Msg.RawEvent.DigitalOptionsPositionChanged1.OrderIds[0]
			tradeData.Type = AssetTypeDigital
			tradeData.ActiveID = res.Msg.ActiveID
			tradeData.TimeFrameInMinutes = max(res.Msg.RawEvent.DigitalOptionsPositionChanged1.InstrumentPeriod/60, 1)
			tradeData.Amount = res.Msg.RawEvent.DigitalOptionsPositionChanged1.BuyAmount
			tradeData.Direction = TradeDirection(res.Msg.RawEvent.DigitalOptionsPositionChanged1.InstrumentDir)
			tradeData.Win = res.Msg.Pnl > 0
			tradeData.OpenTime = res.Msg.OpenTime
			tradeData.Profit = res.Msg.Pnl
		}

		if tradeData.ActiveID == 0 {
			return
		}

		if tradeData.Status == "open" {
			if callbacks.onTradeOpened != nil && lastTradeID[tradeData.Type] != tradeData.OpenTime {
				callbacks.onTradeOpened(tradeData)
				lastTradeID[tradeData.Type] = tradeData.OpenTime
			}
		} else if tradeData.Status == "closed" {
			if callbacks.onTradeClosed != nil {
				callbacks.onTradeClosed(tradeData)
			}
		} else {
			debug.IfVerbose.Printf("Unknown trade status: %s\n", tradeData.Status)
		}
	})
}

func (c *Client) subscribeToTradeChanges() {
	instrumentTypesForSubscription := []string{
		"binary-option",
		"digital-option",
		"turbo-option",
	}

	// for each balance
	for _, balance := range c.balances {

		// and for each instrument type
		for _, instrumentTypeForSubscription := range instrumentTypesForSubscription {
			newRequest := requestEvent{
				Name: "subscribeMessage",
				Msg: map[string]interface{}{
					"name":    "portfolio.position-changed",
					"version": "3.0",
					"params": map[string]interface{}{
						"routingFilters": map[string]interface{}{
							"user_id":         balance.UserID,
							"user_balance_id": balance.ID,
							"instrument_type": instrumentTypeForSubscription,
						},
					},
				},
			}.WithRandomRequestId()

			c.ws.Emit(newRequest)
		}
	}
}
