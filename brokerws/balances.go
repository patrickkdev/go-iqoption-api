package brokerws

import (
	"fmt"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

type BalanceType int

const (
	BalanceTypeReal BalanceType = 1
	BalanceTypeDemo BalanceType = 4
)

type balanceResponse struct {
	RequestID string   `json:"request_id"`
	Name      string   `json:"name"`
	Msg       Balances `json:"msg"`
	Status    int64    `json:"status"`
}

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

func (b *Balances) FindByType(type_ BalanceType) (*Balance, error) {
	for _, balance := range *b {
		if balance.Type == int(type_) {
			return &balance, nil
		}
	}

	return nil, fmt.Errorf("balance not found")
}

func GetBalances(ws *Socket, timeout time.Time) (*Balances, error) {
	eventMsg := map[string]interface{}{
		"name":    "get-balances",
		"body":    map[string]interface{}{},
		"version": "1.0",
	}

	requestEvent := &RequestEvent{
		Name: "sendMessage",
		Msg:  eventMsg,
	}

	resp, err := EmitWithResponse(ws, requestEvent, "balances", timeout)
	if err != nil {
		return nil, err
	}

	responseEvent, err := tjson.Unmarshal[balanceResponse](resp)
	if err != nil {
		return nil, err
	}

	return &responseEvent.Msg, nil
}
