package types

type BinaryTradeData struct {
	Name             string       `json:"name"`
	MicroserviceName string       `json:"microserviceName"`
	Msg              BinaryResult `json:"msg"`
}

type DigitalTradeData struct {
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
			Amount                           float64 `json:"amount"`
			EnrolledAmount                   float64 `json:"enrolled_amount"`
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
	Invest              float64 `json:"invest"`
	InvestEnrolled      float64 `json:"invest_enrolled"`
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
			BuyAmount  float64 `json:"buy_amount"`
			Commission float64 `json:"commission"`
			ExtraData  struct {
				Amount                float64 `json:"amount"`
				Version               string  `json:"version"`
				SpotOption            bool    `json:"spot_option"`
				UseTrailStop          bool    `json:"use_trail_stop"`
				AutoMarginCall        bool    `json:"auto_margin_call"`
				LastChangeReason      string  `json:"last_change_reason"`
				OpenReceivedTime      int     `json:"open_received_time"`
				LowerInstrumentID     string  `json:"lower_instrument_id"`
				UpperInstrumentID     string  `json:"upper_instrument_id"`
				LowerInstrumentStrike int     `json:"lower_instrument_strike"`
				UpperInstrumentStrike int     `json:"upper_instrument_strike"`
				UseTokenForCommission bool    `json:"use_token_for_commission"`
			} `json:"extra_data"`
			LastIndex                 int     `json:"last_index"`
			TpslExtra                 any     `json:"tpsl_extra"`
			SellAmount                float64 `json:"sell_amount"`
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
			BuyAmountEnrolled         float64 `json:"buy_amount_enrolled"`
			CloseEffectAmount         any     `json:"close_effect_amount"`
			CommissionEnrolled        float64 `json:"commission_enrolled"`
			InstrumentActiveID        int     `json:"instrument_active_id"`
			InstrumentIDEscape        string  `json:"instrument_id_escape"`
			SellAmountEnrolled        float64 `json:"sell_amount_enrolled"`
			TakeProfitOrderID         any     `json:"take_profit_order_id"`
			InstrumentExpiration      int     `json:"instrument_expiration"`
			InstrumentUnderlying      string  `json:"instrument_underlying"`
			OpenUnderlyingPrice       float64 `json:"open_underlying_price"`
			PnlRealizedEnrolled       float64 `json:"pnl_realized_enrolled"`
			BuyAvgPriceEnrolled       float64 `json:"buy_avg_price_enrolled"`
			CloseUnderlyingPrice      any     `json:"close_underlying_price"`
			InstrumentStrikeValue     int     `json:"instrument_strike_value"`
			OpenClientPlatformID      int     `json:"open_client_platform_id"`
			SellAvgPriceEnrolled      float64 `json:"sell_avg_price_enrolled"`
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
	Invest                 float64 `json:"invest"`
	InvestEnrolled         float64 `json:"invest_enrolled"`
	SellProfit             float64 `json:"sell_profit"`
	SellProfitEnrolled     float64 `json:"sell_profit_enrolled"`
	ExpectedProfit         float64 `json:"expected_profit"`
	ExpectedProfitEnrolled float64 `json:"expected_profit_enrolled"`
	Pnl                    float64 `json:"pnl"`
	PnlNet                 float64 `json:"pnl_net"`
	CurrentPrice           float64 `json:"current_price"`
	QuoteTimestamp         int     `json:"quote_timestamp"`
	Swap                   int     `json:"swap"`
}

type TradeDigitalResponseEvent struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`
	Msg       struct {
		ID int `json:"id"`
	} `json:"msg"`
	Status int `json:"status"`
}

type TradeBinaryResponseEvent struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`
	Msg       struct {
		UserID             int64       `json:"user_id"`
		ID                 int         `json:"id"`
		RefundValue        int64       `json:"refund_value"`
		Price              float64     `json:"price"`
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
