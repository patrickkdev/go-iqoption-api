package brokerws

import (
	"fmt"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

type binaryCheckResultResponse struct {
	Name             string       `json:"name"`
	MicroserviceName string       `json:"microserviceName"`
	Msg              BinaryResult `json:"msg"`
}

type digitalCheckResultResponse struct {
	Name             string        `json:"name"`
	MicroserviceName string        `json:"microserviceName"`
	Msg              DigitalResult `json:"msg"`
}

type BinaryResult struct {
	RawEvent struct {
		BinaryOptionsOptionChanged1 struct {
			Index                            int64   `json:"index"`
			OptionID                         int64   `json:"option_id"`
			UserID                           int     `json:"user_id"`
			BalanceID                        int     `json:"balance_id"`
			OptionTypeID                     int     `json:"option_type_id"`
			OptionType                       string  `json:"option_type"`
			ActiveID                         int     `json:"active_id"`
			PlatformID                       int     `json:"platform_id"`
			ProfitPercent                    int     `json:"profit_percent"`
			UserBalanceType                  int     `json:"user_balance_type"`
			Currency                         string  `json:"currency"`
			Direction                        string  `json:"direction"`
			Result                           string  `json:"result"`
			Amount                           int     `json:"amount"`
			EnrolledAmount                   int     `json:"enrolled_amount"`
			ProfitAmount                     float64 `json:"profit_amount"`
			WinEnrolledAmount                float64 `json:"win_enrolled_amount"`
			Value                            float64 `json:"value"`
			ExpirationValue                  float64 `json:"expiration_value"`
			OpenTime                         int     `json:"open_time"`
			OpenTimeMillisecond              int64   `json:"open_time_millisecond"`
			ExpirationTime                   int     `json:"expiration_time"`
			ActualExpire                     int     `json:"actual_expire"`
			UserGroupID                      int     `json:"user_group_id"`
			RolloverOptionID                 any     `json:"rollover_option_id"`
			RolloverCommissionOperationID    any     `json:"rollover_commission_operation_id"`
			RolloverCommissionAmount         any     `json:"rollover_commission_amount"`
			RolloverCommissionEnrolledAmount any     `json:"rollover_commission_enrolled_amount"`
			RolloverInitialCommissionAmount  any     `json:"rollover_initial_commission_amount"`
			IsRolledOver                     any     `json:"is_rolled_over"`
			RequestedAt                      int64   `json:"requested_at"`
			CreatedAt                        int64   `json:"created_at"`
			UpdatedAt                        int64   `json:"updated_at"`
		} `json:"binary_options_option_changed1"`
	} `json:"raw_event"`
	Version             int     `json:"version"`
	ID                  string  `json:"id"`
	UserID              int     `json:"user_id"`
	UserBalanceID       int     `json:"user_balance_id"`
	PlatformID          int     `json:"platform_id"`
	ExternalID          int     `json:"external_id"`
	ActiveID            int     `json:"active_id"`
	InstrumentID        string  `json:"instrument_id"`
	Source              string  `json:"source"`
	InstrumentType      string  `json:"instrument_type"`
	Status              string  `json:"status"`
	OpenTime            int     `json:"open_time"`
	OpenQuote           float64 `json:"open_quote"`
	Invest              int     `json:"invest"`
	InvestEnrolled      int     `json:"invest_enrolled"`
	CloseQuote          float64 `json:"close_quote"`
	CloseReason         string  `json:"close_reason"`
	CloseTime           int     `json:"close_time"`
	CloseProfit         float64 `json:"close_profit"`
	CloseProfitEnrolled float64 `json:"close_profit_enrolled"`
	Pnl                 float64 `json:"pnl"`
	PnlRealized         float64 `json:"pnl_realized"`
	PnlNet              float64 `json:"pnl_net"`
	Swap                int     `json:"swap"`
}

type DigitalResult struct {
	RawEvent struct {
		DigitalOptionsPositionChanged1 struct {
			ID         int     `json:"id"`
			Swap       int     `json:"swap"`
			Type       string  `json:"type"`
			Count      float64 `json:"count"`
			Index      int     `json:"index"`
			Status     string  `json:"status"`
			UserID     int     `json:"user_id"`
			CloseAt    any     `json:"close_at"`
			Currency   string  `json:"currency"`
			Leverage   int     `json:"leverage"`
			Pipeline   string  `json:"pipeline"`
			CreateAt   int     `json:"create_at"`
			OrderIds   []int   `json:"order_ids"`
			UpdateAt   int     `json:"update_at"`
			BuyAmount  int     `json:"buy_amount"`
			Commission int     `json:"commission"`
			ExtraData  struct {
				Amount                int    `json:"amount"`
				Version               string `json:"version"`
				SpotOption            bool   `json:"spot_option"`
				UseTrailStop          bool   `json:"use_trail_stop"`
				AutoMarginCall        bool   `json:"auto_margin_call"`
				LastChangeReason      string `json:"last_change_reason"`
				OpenReceivedTime      int    `json:"open_received_time"`
				LowerInstrumentID     string `json:"lower_instrument_id"`
				UpperInstrumentID     string `json:"upper_instrument_id"`
				LowerInstrumentStrike int    `json:"lower_instrument_strike"`
				UpperInstrumentStrike int    `json:"upper_instrument_strike"`
				UseTokenForCommission bool   `json:"use_token_for_commission"`
			} `json:"extra_data"`
			LastIndex                 int     `json:"last_index"`
			TpslExtra                 any     `json:"tpsl_extra"`
			SellAmount                int     `json:"sell_amount"`
			CloseReason               any     `json:"close_reason"`
			PnlRealized               int     `json:"pnl_realized"`
			BuyAvgPrice               float64 `json:"buy_avg_price"`
			CurrencyRate              int     `json:"currency_rate"`
			CurrencyUnit              int     `json:"currency_unit"`
			InstrumentID              string  `json:"instrument_id"`
			SwapEnrolled              int     `json:"swap_enrolled"`
			UserGroupID               int     `json:"user_group_id"`
			CountRealized             int     `json:"count_realized"`
			InstrumentDir             string  `json:"instrument_dir"`
			SellAvgPrice              int     `json:"sell_avg_price"`
			InstrumentType            string  `json:"instrument_type"`
			UserBalanceID             int     `json:"user_balance_id"`
			InstrumentIndex           int     `json:"instrument_index"`
			InstrumentPeriod          int     `json:"instrument_period"`
			InstrumentStrike          float64 `json:"instrument_strike"`
			UserBalanceType           int     `json:"user_balance_type"`
			OpenQuoteTimeMs           int     `json:"open_quote_time_ms"`
			StopLoseOrderID           any     `json:"stop_lose_order_id"`
			BuyAmountEnrolled         int     `json:"buy_amount_enrolled"`
			CloseEffectAmount         any     `json:"close_effect_amount"`
			CommissionEnrolled        int     `json:"commission_enrolled"`
			InstrumentActiveID        int     `json:"instrument_active_id"`
			InstrumentIDEscape        string  `json:"instrument_id_escape"`
			SellAmountEnrolled        int     `json:"sell_amount_enrolled"`
			TakeProfitOrderID         any     `json:"take_profit_order_id"`
			InstrumentExpiration      int     `json:"instrument_expiration"`
			InstrumentUnderlying      string  `json:"instrument_underlying"`
			OpenUnderlyingPrice       float64 `json:"open_underlying_price"`
			PnlRealizedEnrolled       int     `json:"pnl_realized_enrolled"`
			BuyAvgPriceEnrolled       float64 `json:"buy_avg_price_enrolled"`
			CloseUnderlyingPrice      any     `json:"close_underlying_price"`
			InstrumentStrikeValue     int     `json:"instrument_strike_value"`
			OpenClientPlatformID      int     `json:"open_client_platform_id"`
			SellAvgPriceEnrolled      int     `json:"sell_avg_price_enrolled"`
			CloseEffectAmountEnrolled any     `json:"close_effect_amount_enrolled"`
		} `json:"digital_options_position_changed1"`
	} `json:"raw_event"`
	Version                int     `json:"version"`
	ID                     string  `json:"id"`
	UserID                 int     `json:"user_id"`
	UserBalanceID          int     `json:"user_balance_id"`
	PlatformID             int     `json:"platform_id"`
	ExternalID             int     `json:"external_id"`
	ActiveID               int     `json:"active_id"`
	InstrumentID           string  `json:"instrument_id"`
	Source                 string  `json:"source"`
	InstrumentType         string  `json:"instrument_type"`
	Status                 string  `json:"status"`
	OpenTime               int     `json:"open_time"`
	OpenQuote              float64 `json:"open_quote"`
	Invest                 int     `json:"invest"`
	InvestEnrolled         int     `json:"invest_enrolled"`
	SellProfit             float64 `json:"sell_profit"`
	SellProfitEnrolled     float64 `json:"sell_profit_enrolled"`
	ExpectedProfit         int     `json:"expected_profit"`
	ExpectedProfitEnrolled int     `json:"expected_profit_enrolled"`
	Pnl                    float64 `json:"pnl"`
	PnlNet                 float64 `json:"pnl_net"`
	CurrentPrice           float64 `json:"current_price"`
	QuoteTimestamp         int     `json:"quote_timestamp"`
	Swap                   int     `json:"swap"`
}

func CheckResultBinary(ws *Socket, tradeID int, timeout time.Time) (*BinaryResult, bool, error) {
	debug.IfVerbose.Printf("Calling check result binary with tradeID: %d\n", tradeID)

	var err error = nil
	var res binaryCheckResultResponse

	ws.AddEventListener("position-changed", func(event []byte) {
		res, err = tjson.Unmarshal[binaryCheckResultResponse](event)

		debug.IfVerbose.Println("Position changed for binary trade")
		debug.IfVerbose.PrintAsJSON(res)
	})

	for res.Msg.ExternalID != tradeID || res.Msg.Status != "closed" {
		time.Sleep(time.Second)

		if ws.Closed {
			err = fmt.Errorf("websocket closed")
			break
		}

		if time.Since(timeout) > 0 {
			err = fmt.Errorf("timed out waiting for response")
			break
		}

		debug.IfVerbose.Println("Position changed for binary trade:", res.Msg.Status, res.Msg.ExternalID, tradeID)
	}

	ws.RemoveEventListener("position-changed")

	win := res.Msg.CloseReason == "win"
	debug.IfVerbose.Println("Result:", res.Msg.Status, res.Msg.ExternalID, tradeID, win)

	if err != nil {
		return nil, false, err
	}

	return &res.Msg, res.Msg.CloseReason == "win", err
}

func CheckResultDigital(ws *Socket, tradeID int, timeout time.Time) (*DigitalResult, bool, error) {
	debug.IfVerbose.Printf("Calling check result digital with tradeID: %d\n", tradeID)

	var err error = nil
	var res digitalCheckResultResponse

	ws.AddEventListener("position-changed", func(event []byte) {
		res, err = tjson.Unmarshal[digitalCheckResultResponse](event)

		debug.IfVerbose.Println("Position changed for digital trade")
		debug.IfVerbose.PrintAsJSON(res)
	})

	for {
		time.Sleep(time.Second)

		if ws.Closed {
			err = fmt.Errorf("websocket closed")
			break
		}

		if time.Since(timeout) > 0 {
			err = fmt.Errorf("timed out waiting for response")
		}

		orderIDs := res.Msg.RawEvent.DigitalOptionsPositionChanged1.OrderIds

		if len(orderIDs) == 0 {
			continue
		}

		orderID := res.Msg.RawEvent.DigitalOptionsPositionChanged1.OrderIds[0]

		if orderID != tradeID {
			continue
		}

		debug.IfVerbose.Println("Position changed for digital trade:", res.Msg.Status, orderID, tradeID)

		if res.Msg.Status != "closed" {
			continue
		}

		break
	}

	ws.RemoveEventListener("position-changed")

	win := res.Msg.Pnl > 0
	debug.IfVerbose.Println("Check result digital: ", res.Msg.Pnl, win)

	if err != nil {
		return nil, false, err
	}

	return &res.Msg, win, err
}

func GetSubscriptionsToPositionChangedEvent(userID int, userBalanceID int) []*RequestEvent {
	var intrumentTypesForSubscription = []string{
		"binary-option",
		"digital-option",
		"turbo-option",
	}

	var requestEvents = make([]*RequestEvent, 0)

	for _, instrumentTypeForSubscription := range intrumentTypesForSubscription {
		newRequest := &RequestEvent{
			Name: "subscribeMessage",
			Msg: map[string]interface{}{
				"name":    "portfolio.position-changed",
				"version": "3.0",
				"params": map[string]interface{}{
					"routingFilters": map[string]interface{}{
						"user_id":         userID,
						"user_balance_id": userBalanceID,
						"instrument_type": instrumentTypeForSubscription,
					},
				},
			},
		}

		requestEvents = append(requestEvents, newRequest)
	}

	return requestEvents
}
