package brokerws

import (
	"strings"

	"github.com/patrickkdev/Go-IQOption-API/btypes"
	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

type TradeChangedCallbacks struct {
	OnTradeOpenedCallback func(tradeData btypes.TradeData)
	OnTradeClosedCallback func(tradeData btypes.TradeData)
}

func (ws *Socket) OnTradeChanged(callbacks TradeChangedCallbacks) {
	lastTradeID := make(map[btypes.AssetType]int)

	ws.AddEventListener("position-changed", func(event []byte) {
		eventString := string(event)
		tradeData := btypes.TradeData{ActiveID: 0}

		if strings.Contains(eventString, "binary_options_option_changed1") {
			res, err := tjson.Unmarshal[binaryTradeData](event)
			if err != nil {
				return
			}

			tradeData.Status = res.Msg.Status
			tradeData.TradeID = res.Msg.ExternalID
			tradeData.Type = btypes.AssetTypeBinary
			tradeData.ActiveID = res.Msg.ActiveID
			tradeData.TimeFrameInMinutes = max((res.Msg.RawEvent.BinaryOptionsOptionChanged1.ExpirationTime-res.Msg.RawEvent.BinaryOptionsOptionChanged1.OpenTime)/60, 1)
			tradeData.Amount = res.Msg.RawEvent.BinaryOptionsOptionChanged1.Amount
			tradeData.Direction = btypes.TradeDirection(res.Msg.RawEvent.BinaryOptionsOptionChanged1.Direction)
			tradeData.Win = res.Msg.CloseReason == "win"
			tradeData.OpenTime = res.Msg.OpenTime
			tradeData.Profit = res.Msg.Pnl

		} else if strings.Contains(eventString, "digital_options_position_changed1") {
			res, err := tjson.Unmarshal[digitalTradeData](event)
			if err != nil {
				return
			}

			tradeData.Status = res.Msg.Status
			tradeData.TradeID = res.Msg.RawEvent.DigitalOptionsPositionChanged1.OrderIds[0]
			tradeData.Type = btypes.AssetTypeDigital
			tradeData.ActiveID = res.Msg.ActiveID
			tradeData.TimeFrameInMinutes = max(res.Msg.RawEvent.DigitalOptionsPositionChanged1.InstrumentPeriod/60, 1)
			tradeData.Amount = res.Msg.RawEvent.DigitalOptionsPositionChanged1.BuyAmount
			tradeData.Direction = btypes.TradeDirection(res.Msg.RawEvent.DigitalOptionsPositionChanged1.InstrumentDir)
			tradeData.Win = res.Msg.Pnl > 0
			tradeData.OpenTime = res.Msg.OpenTime
			tradeData.Profit = res.Msg.Pnl
		}

		if tradeData.ActiveID == 0 {
			return
		}

		if tradeData.Status == "open" {
			if callbacks.OnTradeOpenedCallback != nil && lastTradeID[tradeData.Type] != tradeData.OpenTime {
				callbacks.OnTradeOpenedCallback(tradeData)
				lastTradeID[tradeData.Type] = tradeData.OpenTime
			}
		} else if tradeData.Status == "closed" {
			if callbacks.OnTradeClosedCallback != nil {
				callbacks.OnTradeClosedCallback(tradeData)
			}
		} else {
			debug.IfVerbose.Printf("Unknown trade status: %s\n", tradeData.Status)
		}
	})
}

func (ws *Socket) SubscribeToTradeChanges(balances btypes.Balances) {
	instrumentTypesForSubscription := []string{
		"binary-option",
		"digital-option",
		"turbo-option",
	}

	// for each balance
	for _, balance := range balances {

		// and for each instrument type
		for _, instrumentTypeForSubscription := range instrumentTypesForSubscription {
			newRequest := &btypes.RequestEvent{
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
			}

			ws.Emit(newRequest)
		}
	}

}
