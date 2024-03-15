package wsapi

// type binaryCheckResultResponse struct {
// 	Name             string       `json:"name"`
// 	MicroserviceName string       `json:"microserviceName"`
// 	Msg              BinaryResult `json:"msg"`
// }

// type BinaryResult struct {
// 	RawEvent struct {
// 		BinaryOptionsOptionChanged1 struct {
// 			Index                            int64   `json:"index"`
// 			OptionID                         int64   `json:"option_id"`
// 			UserID                           int     `json:"user_id"`
// 			BalanceID                        int     `json:"balance_id"`
// 			OptionTypeID                     int     `json:"option_type_id"`
// 			OptionType                       string  `json:"option_type"`
// 			ActiveID                         int     `json:"active_id"`
// 			PlatformID                       int     `json:"platform_id"`
// 			ProfitPercent                    int     `json:"profit_percent"`
// 			UserBalanceType                  int     `json:"user_balance_type"`
// 			Currency                         string  `json:"currency"`
// 			Direction                        string  `json:"direction"`
// 			Result                           string  `json:"result"`
// 			Amount                           int     `json:"amount"`
// 			EnrolledAmount                   int     `json:"enrolled_amount"`
// 			ProfitAmount                     float64 `json:"profit_amount"`
// 			WinEnrolledAmount                float64 `json:"win_enrolled_amount"`
// 			Value                            float64 `json:"value"`
// 			ExpirationValue                  float64 `json:"expiration_value"`
// 			OpenTime                         int     `json:"open_time"`
// 			OpenTimeMillisecond              int64   `json:"open_time_millisecond"`
// 			ExpirationTime                   int     `json:"expiration_time"`
// 			ActualExpire                     int     `json:"actual_expire"`
// 			UserGroupID                      int     `json:"user_group_id"`
// 			RolloverOptionID                 any     `json:"rollover_option_id"`
// 			RolloverCommissionOperationID    any     `json:"rollover_commission_operation_id"`
// 			RolloverCommissionAmount         any     `json:"rollover_commission_amount"`
// 			RolloverCommissionEnrolledAmount any     `json:"rollover_commission_enrolled_amount"`
// 			RolloverInitialCommissionAmount  any     `json:"rollover_initial_commission_amount"`
// 			IsRolledOver                     any     `json:"is_rolled_over"`
// 			RequestedAt                      int64   `json:"requested_at"`
// 			CreatedAt                        int64   `json:"created_at"`
// 			UpdatedAt                        int64   `json:"updated_at"`
// 		} `json:"binary_options_option_changed1"`
// 	} `json:"raw_event"`
// 	Version             int64   `json:"version"`
// 	ID                  string  `json:"id"`
// 	UserID              int     `json:"user_id"`
// 	UserBalanceID       int     `json:"user_balance_id"`
// 	PlatformID          int     `json:"platform_id"`
// 	ExternalID          int64   `json:"external_id"`
// 	ActiveID            int     `json:"active_id"`
// 	InstrumentID        string  `json:"instrument_id"`
// 	Source              string  `json:"source"`
// 	InstrumentType      string  `json:"instrument_type"`
// 	Status              string  `json:"status"`
// 	OpenTime            int64   `json:"open_time"`
// 	OpenQuote           float64 `json:"open_quote"`
// 	Invest              int     `json:"invest"`
// 	InvestEnrolled      int     `json:"invest_enrolled"`
// 	CloseQuote          float64 `json:"close_quote"`
// 	CloseReason         string  `json:"close_reason"`
// 	CloseTime           int64   `json:"close_time"`
// 	CloseProfit         float64 `json:"close_profit"`
// 	CloseProfitEnrolled float64 `json:"close_profit_enrolled"`
// 	Pnl                 float64 `json:"pnl"`
// 	PnlRealized         float64 `json:"pnl_realized"`
// 	PnlNet              float64 `json:"pnl_net"`
// 	Swap                int     `json:"swap"`
// }

// func CheckResultBinary (ws *Socket, tradeID int) (result *binaryCheckResultResponse, err error) {

// }

func GetSubscriptionToPositionChanged(userID int, userBalanceID int) *Event {
	return &Event{
		Name: "subscribeMessage",
		Msg: map[string]interface{}{
			"name":    "portfolio.position-changed",
			"version": "3.0",
			"params": map[string]interface{}{
				"routingFilters": map[string]interface{}{
					"user_id":         userID,
					"user_balance_id": userBalanceID,
					"instrument_type": "exchange-option",
				},
			},
		},
	}
}
