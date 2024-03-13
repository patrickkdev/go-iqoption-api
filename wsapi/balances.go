package wsapi

import (
	"fmt"
	"patrickkdev/Go-IQOption-API/tjson"
	"time"
)

type balanceResponse struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`      
	Msg       Balances  `json:"msg"`       
	Status    int64  `json:"status"`    
}

type Balance struct {
	ID                int64       `json:"id"`                 
	UserID            int64       `json:"user_id"`            
	Type              int64       `json:"type"`               
	Amount            float64     `json:"amount"`             
	EnrolledAmount    float64     `json:"enrolled_amount"`    
	EnrolledSumAmount float64     `json:"enrolled_sum_amount"`
	BonusAmount       int64       `json:"bonus_amount"`       
	HoldAmount        int64       `json:"hold_amount"`        
	OrdersAmount      int64       `json:"orders_amount"`      
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

func (b *Balances) FindByType(type_ int64) (*Balance, error) {
	for _, balance := range *b {
		if balance.Type == type_ {
			return &balance, nil
		}
	}

	return nil, fmt.Errorf("balance not found")
}

func GetBalances(ws *Socket, serverTimeStamp int) (*Balances, error) {
	eventMsg := map[string]interface{}{
		"name": "get-balances",
		"body": map[string]interface{}{},
		"version": "1.0",
	}

	requestEvent := &Event{
		Name:      "sendMessage",
		Msg:       eventMsg,
		RequestId: fmt.Sprint(serverTimeStamp),
	}

	resp, err := EmitWithResponse(ws, requestEvent, "balances", time.Now().Add(10*time.Second))
	if err != nil {
		return nil, err
	}

	responseEvent, err := tjson.Unmarshal[balanceResponse](resp)
	if err != nil {
		return nil, err
	}

	return &responseEvent.Msg, nil
}