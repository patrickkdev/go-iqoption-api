package wsapi

import "fmt"

type TradeDirection string

const (
	TradeDirectionCall TradeDirection = "call"
	TradeDirectionPut  TradeDirection = "put"
)

type TradeType string

const (
	TradeTypeDigital TradeType = "digital"
	TradeTypeBinary  TradeType = "binary"
)

type TradeBalance int

const (
	TradeBalanceDemo TradeBalance = 4
)

func TradeDigital(ws *Socket, amount float64, direction TradeDirection, activeID int, duration int, targetBalanceID int, serverTimeStamp int64) (orderID int, err error) {
	eventMsg := map[string]interface{}{
		"name":    "digital-options.place-digital-option",
		"version": "3.0",
		"body": map[string]interface{}{
			"amount":           amount,
			"asset_id":         activeID,
			"instrument_id":    "do76A20240312D204000T1MCSPT",
			"instrument_index": 2307938,
			"user_balance_id":  targetBalanceID,
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

func TradeBinary(ws *Socket, amount float64, direction TradeDirection, activeID int, duration int, targetBalanceID int, serverTimeStamp int64) {
	exp, idx := GetExpirationTime(serverTimeStamp, duration)

	fmt.Println("expiration time: ", idx)

	optionTypeID := map[bool]int{
		true:  3, // turbo
		false: 1, // binary
	}[idx < 5]

	eventMsg := map[string]interface{}{
		"name":    "binary-options.open-option",
		"version": "1.0",
		"body": map[string]interface{}{
			"price":           amount,
			"active_id":       activeID,
			"expired":         exp,
			"direction":       direction,
			"option_type_id":  optionTypeID,
			"user_balance_id": targetBalanceID,
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
