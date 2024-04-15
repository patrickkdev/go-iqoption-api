package broker

import (
	"encoding/json"
	"fmt"

	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
	"github.com/patrickkdev/Go-IQOption-API/internal/types"
)

type BalanceType int

const (
	BalanceTypeReal BalanceType = 1
	BalanceTypeDemo BalanceType = 4
)

type BalanceResponse struct {
	RequestID string    `json:"request_id"`
	Name      string    `json:"name"`
	Msg       []Balance `json:"msg"`
	Status    int64     `json:"status"`
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

// Although balances are kept updated by subscription to 'balance-changed' event,
// this function can be used to force an update and return the latest balances
func (c *Client) GetUpdatedBalances() (Balances, error) {
	requestEvent := requestEvent{
		Name: "sendMessage",
		Msg: map[string]interface{}{
			"name":    "get-balances",
			"body":    map[string]interface{}{},
			"version": "1.0",
		},
	}.WithRandomRequestId()

	resp, err := c.ws.EmitWithResponse(requestEvent, "balances", c.getTimeout())
	if err != nil {
		return nil, err
	}

	responseEvent, err := tjson.Unmarshal[BalanceResponse](resp)
	if err != nil {
		return nil, err
	}

	c.balances = responseEvent.Msg

	return responseEvent.Msg, nil
}

func (c *Client) keepBalancesUpdated() {
	requestEvent := requestEvent{
		Name: "subscribeMessage",
		Msg: map[string]interface{}{
			"name":    "internal-billing.balance-changed",
			"version": "1.0",
		},
	}.WithRandomRequestId()

	c.ws.Subscribe(requestEvent, "balance-changed", func(event []byte) {
		var update types.BalanceUpdateEvent
		json.Unmarshal(event, &update)

		c.balances.update(update.Msg.CurrentBalance)
	})
}

func (b *Balances) FindByType(type_ BalanceType) (Balance, error) {
	for _, balance := range []Balance(*b) {
		if balance.Type == int(type_) {
			return balance, nil
		}
	}

	return Balance{}, fmt.Errorf("balance not found")
}

func (b *Balances) update(update types.BalanceUpdate) {
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
