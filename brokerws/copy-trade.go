package brokerws

import (
	"strings"

	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

type TradeData struct {
	Type              AssetType
	Direction          TradeDirection
	TimeFrameInMinutes int
	ActiveID           int
	Amount             float64
}

func OnOpenTrade(ws *Socket, callback func(tradeData TradeData)) {
	ws.AddEventListener("position-changed", func(event []byte) {
		eventString := string(event)

		tradeData := TradeData{ActiveID: 0}
		if strings.Contains(eventString, "binary_options_option_changed1") {
			res, err := tjson.Unmarshal[binaryTradeData](event)
			if err != nil || res.Msg.Status != "open" {
				return
			}

			tradeData.Type =              AssetTypeBinary
			tradeData.ActiveID =           res.Msg.ActiveID
			tradeData.TimeFrameInMinutes = (res.Msg.RawEvent.BinaryOptionsOptionChanged1.ExpirationTime - res.Msg.RawEvent.BinaryOptionsOptionChanged1.OpenTime) / 60
			tradeData.Amount =             res.Msg.RawEvent.BinaryOptionsOptionChanged1.Amount
			tradeData.Direction =          TradeDirection(res.Msg.RawEvent.BinaryOptionsOptionChanged1.Direction)
			
		} else if strings.Contains(eventString, "digital_options_position_changed1") {
			res, err := tjson.Unmarshal[digitalTradeData](event)
			if err != nil || res.Msg.Status != "open" {
				return
			}

			tradeData.Type =              AssetTypeDigital
			tradeData.ActiveID =           res.Msg.ActiveID
			tradeData.TimeFrameInMinutes = res.Msg.RawEvent.DigitalOptionsPositionChanged1.InstrumentPeriod / 60
			tradeData.Amount =             res.Msg.RawEvent.DigitalOptionsPositionChanged1.BuyAmount
			tradeData.Direction =          TradeDirection(res.Msg.RawEvent.DigitalOptionsPositionChanged1.InstrumentDir)
		}

		if tradeData.ActiveID == 0 {
			return
		}

		callback(tradeData)
	})
}