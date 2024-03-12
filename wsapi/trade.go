package wsapi

import "fmt"

func TradeDigital(ws *Socket, amount float64, assetID int, userBalanceID int, serverTimeStamp int) (orderID int, err error) {
	eventMsg := map[string]interface{}{
		"name":    "digital-options.place-digital-option",
		"version": "3.0",
		"body": map[string]interface{}{
			"amount":           amount,
			"asset_id":         assetID,
			"instrument_id":    "do76A20240312D204000T1MCSPT",
			"instrument_index": 2307938,
			"user_balance_id":  userBalanceID,
		},
	}

	requestEvent := &Event{
		Name:      "sendMessage",
		Msg:       eventMsg,
		RequestId: fmt.Sprint(serverTimeStamp),
		LocalTime: serverTimeStamp,
	}

	ws.EmitEvent(requestEvent)
	
	return 2, nil
}