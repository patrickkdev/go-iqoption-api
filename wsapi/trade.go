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

func TradeBinary(ws *Socket, amount float64, activeID int, userBalanceID int, serverTimeStamp int) {
	eventMsg := map[string]interface{}{
		"name":    "binary-options.open-option",
		"version": "2.0",
		"body": map[string]interface{}{
			"user_balance_id":  	userBalanceID,
			"active_id":         	activeID,
			"option_type_id":    	3,
			"direction":         	"put",
			"expired":           	serverTimeStamp + 60,
			"refund_value":      	0,
			"amount":           	amount,
			"value":            	0,
			"profit_percent":    	85,
		},
	}

	requestEvent := &Event{
		Name:      "sendMessage",
		Msg:       eventMsg,
		RequestId: fmt.Sprint(serverTimeStamp),
		LocalTime: serverTimeStamp,
	}

	ws.EmitEvent(requestEvent)
}