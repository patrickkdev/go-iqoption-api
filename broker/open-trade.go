package broker

import (
	"fmt"

	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
	"github.com/patrickkdev/Go-IQOption-API/internal/types"
)

// Immediately opens a new trade if params are valid and the asset is available
func (c *Client) OpenTrade(type_ AssetType, amount float64, direction TradeDirection, activeID int, timeFrameInMinutes int, balance BalanceType) (int, error) {
	targetBalance, err := c.balances.FindByType(balance)
	if err != nil {
		return 0, fmt.Errorf("invalid balance type: %d. %s", balance, err.Error())
	}

	switch type_ {
	case AssetTypeBinary:
		return c.openBinaryTrade(
			amount,
			direction,
			activeID,
			timeFrameInMinutes,
			targetBalance.ID,
		)
	case AssetTypeDigital:
		return c.openDigitalTrade(
			amount,
			direction,
			activeID,
			timeFrameInMinutes,
			targetBalance.ID,
		)
	}

	return 0, nil
}

func (c *Client) openDigitalTrade(amount float64, direction TradeDirection, activeID int, duration int, targetBalanceID int) (int, error) {
	exp, _ := getExpirationTime(c.serverTimestamp, duration)

	instrument, err := c.getInstrument(activeID, exp)
	if err != nil {
		return 0, err
	}

	instrumentID, err := instrument.Data.GetSymbol(string(direction))
	if err != nil {
		return 0, err
	}

	requestEvent := requestEvent{
		Name: "sendMessage",
		Msg: map[string]interface{}{
			"name":    "digital-options.place-digital-option",
			"version": "3.0",
			"body": map[string]interface{}{
				"amount":           fmt.Sprint(amount),
				"asset_id":         activeID,
				"instrument_id":    instrumentID,
				"instrument_index": instrument.Index,
				"user_balance_id":  targetBalanceID,
			},
		},
	}.WithRandomRequestId()

	resp, err := c.ws.EmitWithResponse(requestEvent, "digital-option-placed", c.getTimeout())
	if err != nil {
		return 0, err
	}

	responseEvent, err := tjson.Unmarshal[types.TradeDigitalResponseEvent](resp)
	if err != nil {
		return 0, err
	}

	debug.IfVerbose.PrintAsJSON(map[string]any{"responseEvent": responseEvent})

	if responseEvent.Msg.ID == 0 {
		return 0, fmt.Errorf("error placing trade")
	}

	return responseEvent.Msg.ID, nil
}

func (c *Client) openBinaryTrade(amount float64, direction TradeDirection, activeID int, duration int, targetBalanceID int) (int, error) {
	exp, idx := getExpirationTime(c.serverTimestamp, duration)

	optionTypeID := map[bool]int{
		true:  map[bool]int{true: 12, false: 3}[c.GetBrokerDomain() == "trade.bull-ex.com"], // turbo
		false: 1,                                                                            // binary
	}[idx < 5]

	requestEvent := requestEvent{
		Name: "sendMessage",
		Msg: map[string]interface{}{
			"name":    "binary-options.open-option",
			"version": "1.0",
			"body": map[string]interface{}{
				"price":           amount,
				"active_id":       activeID,
				"expired":         exp,
				"direction":       direction,
				"option_type_id":  optionTypeID,
				"expiration_size": duration * 60,
				"user_balance_id": targetBalanceID,
			},
		},
	}.WithRandomRequestId()

	resp, err := c.ws.EmitWithResponse(requestEvent, "option", c.getTimeout())
	if err != nil {
		return 0, err
	}

	responseEvent, err := tjson.Unmarshal[types.TradeBinaryResponseEvent](resp)
	if err != nil {
		return 0, err
	}

	debug.IfVerbose.PrintAsJSON(map[string]any{"responseEvent": responseEvent})

	if responseEvent.Msg.ID == 0 {
		return 0, fmt.Errorf("error placing trade")
	}

	return responseEvent.Msg.ID, nil
}
