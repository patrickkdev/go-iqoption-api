package wsapi

import (
	"fmt"
	"patrickkdev/Go-IQOption-API/tjson"
	"time"
)

type instrumentsResponseEvent struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`
	Msg       struct {
		Instruments []instrument `json:"instruments"`
		NotFound []any `json:"not_found"`
	} `json:"msg"`
	Status int `json:"status"`
}

type instrument struct {
	Index          	int     `json:"index"`
	InstrumentType 	string  `json:"instrument_type"`
	AssetID        	int     `json:"asset_id"`
	TradingGroupID 	string  `json:"trading_group_id"`
	Expiration     	int     `json:"expiration"`
	Period         	int     `json:"period"`
	Quote          	float64 `json:"quote"`
	Data           	instrumentData `json:"data"`
	Volatility     	float64 `json:"volatility"`
	GeneratedAt    	int     `json:"generated_at"`
	Deadtime        int `json:"deadtime"`
	BuybackDeadtime int `json:"buyback_deadtime"`
}

type instrumentData []struct {
		Strike    string    `json:"strike"`   
		Symbol    string    `json:"symbol"`   
		Direction TradeDirection `json:"direction"`
}

func (data instrumentData) GetSymbol(direction TradeDirection) (symbol string, err error) {
	for _, v := range data {
		if v.Strike == "SPT" && v.Direction == direction {
			return v.Symbol, nil
		}
	}

	return "", fmt.Errorf("no symbol found for direction %s", direction)
}

func GetInstrument(ws *Socket, activeID int, exp int, serverTimestamp int64) (*instrument, error) {
	msg := map[string]interface{}{
		"name":       "digital-option-instruments.get-instruments",
		"version":    "2.0",
		"body":    		 map[string]interface{}{
			"instrument_type": "digital-option",
			"asset_id": activeID,
    },
	}

	requestEvent := &RequestEvent{
		Name:      "sendMessage",
		Msg:       msg,
		RequestId: fmt.Sprint(serverTimestamp),
	}

	resp, err := EmitWithResponse(ws, requestEvent, "instruments", time.Now().Add(1*time.Minute))
	if err != nil {
		return nil, err
	}

	responseEvent, err := tjson.Unmarshal[instrumentsResponseEvent](resp)
	if err != nil {
		return nil, err
	}

	if len(responseEvent.Msg.Instruments) == 0 {
		return nil, fmt.Errorf("no instrument found for active id %d", activeID)
	}

	for _, instrument := range responseEvent.Msg.Instruments {
		if instrument.AssetID == activeID && instrument.Expiration == exp {
			return &instrument, nil
		}
	}

	return nil, fmt.Errorf("no instrument found for active id %d and expiration %d", activeID, exp)
}