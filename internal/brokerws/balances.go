package brokerws

import (
	"encoding/json"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/btypes"
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

type balanceResponse struct {
	RequestID string          `json:"request_id"`
	Name      string          `json:"name"`
	Msg       btypes.Balances `json:"msg"`
	Status    int64           `json:"status"`
}

func (ws *Socket) GetBalances(timeout time.Time) (btypes.Balances, error) {
	eventMsg := map[string]interface{}{
		"name":    "get-balances",
		"body":    map[string]interface{}{},
		"version": "1.0",
	}

	requestEvent := &btypes.RequestEvent{
		Name: "sendMessage",
		Msg:  eventMsg,
	}

	resp, err := ws.EmitWithResponse(requestEvent, "balances", timeout)
	if err != nil {
		return nil, err
	}

	responseEvent, err := tjson.Unmarshal[balanceResponse](resp)
	if err != nil {
		return nil, err
	}

	return responseEvent.Msg, nil
}

type balanceUpdateEvent struct {
	Name             string `json:"name"`
	MicroserviceName string `json:"microserviceName"`
	Msg              struct {
		CurrentBalance btypes.BalanceUpdate `json:"current_balance"`
		ID             int                  `json:"id"`
		UserID         int                  `json:"user_id"`
	} `json:"msg"`
}

func (ws *Socket) OnBalanceChanged(callback func(update btypes.BalanceUpdate)) {
	requestEvent := &btypes.RequestEvent{
		Name: "subscribeMessage",
		Msg: map[string]interface{}{
			"name":    "internal-billing.balance-changed",
			"version": "1.0",
		},
	}

	ws.Emit(requestEvent)

	ws.AddEventListener("balance-changed", func(event []byte) {
		var update balanceUpdateEvent
		json.Unmarshal(event, &update)

		callback(update.Msg.CurrentBalance)
	})
}
