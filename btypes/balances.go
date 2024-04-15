package btypes

import "fmt"

type BalanceType int

const (
	BalanceTypeReal BalanceType = 1
	BalanceTypeDemo BalanceType = 4
)

type Balance struct {
	ID                int         `json:"id"`
	UserID            int         `json:"user_id"`
	Type              int         `json:"type"`
	Amount            float64     `json:"amount"`
	EnrolledAmount    float64     `json:"enrolled_amount"`
	EnrolledSumAmount float64     `json:"enrolled_sum_amount"`
	BonusAmount       int         `json:"bonus_amount"`
	HoldAmount        int         `json:"hold_amount"`
	OrdersAmount      int         `json:"orders_amount"`
	Currency          string      `json:"currency"`
	TournamentID      interface{} `json:"tournament_id"`
	TournamentName    interface{} `json:"tournament_name"`
	IsFiat            bool        `json:"is_fiat"`
	IsMarginal        bool        `json:"is_marginal"`
	HasDeposits       bool        `json:"has_deposits"`
	AuthAmount        int64       `json:"auth_amount"`
	Equivalent        int64       `json:"equivalent"`
}

type Balances []Balance

func (b *Balances) FindByType(type_ BalanceType) (Balance, error) {
	for _, balance := range []Balance(*b) {
		if balance.Type == int(type_) {
			return balance, nil
		}
	}

	return Balance{}, fmt.Errorf("balance not found")
}

func (b *Balances) Update(update BalanceUpdate) {
	for i, balance := range *b {
		if balance.ID == update.ID {
			(*b)[i].Amount = update.NewAmount
			(*b)[i].EnrolledAmount = update.EnrolledAmount
			(*b)[i].BonusAmount = update.BonusAmount
			(*b)[i].IsFiat = update.IsFiat
			(*b)[i].IsMarginal = update.IsMarginal
			return
		}
	}
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
