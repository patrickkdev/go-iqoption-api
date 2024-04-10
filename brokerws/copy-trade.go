package brokerws

import (
	"strings"

	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

type TradeData struct {
	type_              AssetType
	activeID           int
	timeFrameInMinutes int
	amount             float64
	direction          TradeDirection
}

func OnOpenTrade(ws *Socket, callback func(tradeData TradeData)) {
	ws.AddEventListener("position-changed", func(event []byte) {
		eventString := string(event)

		tradeData := TradeData{activeID: 0}

		if strings.Contains(eventString, "binary_options_option_changed1") {
			res, err := tjson.Unmarshal[binaryTradeData](event)
			if err != nil || res.Msg.Status != "open" {
				return
			}

			tradeData.type_ =              AssetTypeBinary
			tradeData.activeID =           res.Msg.ActiveID
			tradeData.timeFrameInMinutes = (res.Msg.RawEvent.BinaryOptionsOptionChanged1.ExpirationTime - res.Msg.RawEvent.BinaryOptionsOptionChanged1.OpenTime) / 60
			tradeData.amount =             res.Msg.RawEvent.BinaryOptionsOptionChanged1.Amount
			tradeData.direction =          TradeDirection(res.Msg.RawEvent.BinaryOptionsOptionChanged1.Direction)
			
		} else if strings.Contains(eventString, "digital_options_position_changed1") {
			res, err := tjson.Unmarshal[digitalTradeData](event)
			if err != nil || res.Msg.Status != "open" {
				return
			}

			tradeData.type_ =              AssetTypeDigital
			tradeData.activeID =           res.Msg.ActiveID
			tradeData.timeFrameInMinutes = res.Msg.RawEvent.DigitalOptionsPositionChanged1.InstrumentPeriod / 60
			tradeData.amount =             res.Msg.RawEvent.DigitalOptionsPositionChanged1.BuyAmount
			tradeData.direction =          TradeDirection(res.Msg.RawEvent.DigitalOptionsPositionChanged1.InstrumentDir)
		}

		if tradeData.activeID == 0 {
			return
		}

		callback(tradeData)
	})
}