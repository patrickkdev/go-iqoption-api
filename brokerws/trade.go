package brokerws

import (
	"fmt"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

type TradeDirection string

const (
	TradeDirectionCall TradeDirection = "call"
	TradeDirectionPut  TradeDirection = "put"
)

type TradeShouldWaitForResult bool

const (
	WaitForResult      TradeShouldWaitForResult = true
	DoNotWaitForResult TradeShouldWaitForResult = false
)

type tradeDigitalResponseEvent struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`
	Msg       struct {
		ID int `json:"id"`
	} `json:"msg"`
	Status int `json:"status"`
}

type tradeBinaryResponseEvent struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`
	Msg       struct {
		UserID             int64       `json:"user_id"`
		ID                 int         `json:"id"`
		RefundValue        int64       `json:"refund_value"`
		Price              int64       `json:"price"`
		Exp                int64       `json:"exp"`
		Created            int64       `json:"created"`
		CreatedMillisecond int64       `json:"created_millisecond"`
		TimeRate           int64       `json:"time_rate"`
		Type               string      `json:"type"`
		Act                int64       `json:"act"`
		Direction          string      `json:"direction"`
		ExpValue           int64       `json:"exp_value"`
		Value              float64     `json:"value"`
		ProfitIncome       int64       `json:"profit_income"`
		ProfitReturn       int64       `json:"profit_return"`
		RobotID            interface{} `json:"robot_id"`
		ClientPlatformID   int64       `json:"client_platform_id"`
	} `json:"msg"`
	Status int64 `json:"status"`
}

func TradeDigital(ws *Socket, amount float64, direction TradeDirection, activeID int, duration int, targetBalanceID int, serverTimeStamp int64, timeout time.Time) (tradeID int, err error) {
	exp, _ := GetExpirationTime(serverTimeStamp, duration)

	instrument, err := GetInstrument(ws, activeID, exp, timeout)
	if err != nil {
		return 0, err
	}

	instrumentID, err := instrument.Data.GetSymbol(direction)
	if err != nil {
		return 0, err
	}

	eventMsg := map[string]interface{}{
		"name":    "digital-options.place-digital-option",
		"version": "3.0",
		"body": map[string]interface{}{
			"amount":           fmt.Sprint(amount),
			"asset_id":         activeID,
			"instrument_id":    instrumentID,
			"instrument_index": instrument.Index,
			"user_balance_id":  targetBalanceID,
		},
	}

	requestEvent := &RequestEvent{
		Name: "sendMessage",
		Msg:  eventMsg,
	}

	debug.IfVerbose.Println("requestEvent:")
	debug.IfVerbose.PrintAsJSON(requestEvent)

	resp, err := EmitWithResponse(ws, requestEvent, "digital-option-placed", timeout)
	if err != nil {
		return 0, err
	}

	responseEvent, err := tjson.Unmarshal[tradeDigitalResponseEvent](resp)
	if err != nil {
		return 0, err
	}

	debug.IfVerbose.Println("responseEvent:")
	debug.IfVerbose.PrintAsJSON(responseEvent)

	if responseEvent.Msg.ID == 0 {
		return 0, fmt.Errorf("error placing trade")
	}

	return responseEvent.Msg.ID, nil
}

func TradeBinary(ws *Socket, amount float64, direction TradeDirection, activeID int, duration int, targetBalanceID int, serverTimeStamp int64, timeout time.Time) (tradeID int, err error) {
	exp, idx := GetExpirationTime(serverTimeStamp, duration)
	debug.IfVerbose.Println("expiration time: ", idx)

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

	requestEvent := &RequestEvent{
		Name: "sendMessage",
		Msg:  eventMsg,
	}

	debug.IfVerbose.Println("requestEvent:")
	debug.IfVerbose.PrintAsJSON(requestEvent)

	resp, err := EmitWithResponse(ws, requestEvent, "option", timeout)
	if err != nil {
		return 0, err
	}

	responseEvent, err := tjson.Unmarshal[tradeBinaryResponseEvent](resp)
	if err != nil {
		return 0, err
	}

	debug.IfVerbose.Println("responseEvent:")
	debug.IfVerbose.PrintAsJSON(responseEvent)

	if responseEvent.Msg.ID == 0 {
		return 0, fmt.Errorf("error placing trade")
	}

	return responseEvent.Msg.ID, nil
}
