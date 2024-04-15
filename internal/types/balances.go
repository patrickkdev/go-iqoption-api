package types

type BalanceUpdateEvent struct {
	Name             string `json:"name"`
	MicroserviceName string `json:"microserviceName"`
	Msg              struct {
		CurrentBalance BalanceUpdate `json:"current_balance"`
		ID             int           `json:"id"`
		UserID         int           `json:"user_id"`
	} `json:"msg"`
}

type BalanceUpdate struct {
	ID                  int     `json:"id"`
	Amount              float64 `json:"amount"`
	EnrolledAmount      float64 `json:"enrolled_amount"`
	BonusAmount         int     `json:"bonus_amount"`
	BonusEnrolledAmount int     `json:"bonus_enrolled_amount"`
	Currency            string  `json:"currency"`
	Type                int     `json:"type"`
	Index               int64   `json:"index"`
	IsFiat              bool    `json:"is_fiat"`
	NewAmount           float64 `json:"new_amount"`
	BonusTotalAmount    int     `json:"bonus_total_amount"`
	IsMarginal          bool    `json:"is_marginal"`
}
