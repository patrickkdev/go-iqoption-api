package types

import "fmt"

type InstrumentsResponseEvent struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`
	Msg       struct {
		Instruments []Instrument `json:"instruments"`
		NotFound    []any        `json:"not_found"`
	} `json:"msg"`
	Status int `json:"status"`
}

type Instrument struct {
	Index           int            `json:"index"`
	InstrumentType  string         `json:"instrument_type"`
	AssetID         int            `json:"asset_id"`
	TradingGroupID  string         `json:"trading_group_id"`
	Expiration      int            `json:"expiration"`
	Period          int            `json:"period"`
	Quote           float64        `json:"quote"`
	Data            InstrumentData `json:"data"`
	Volatility      float64        `json:"volatility"`
	GeneratedAt     int            `json:"generated_at"`
	Deadtime        int            `json:"deadtime"`
	BuybackDeadtime int            `json:"buyback_deadtime"`
}

type InstrumentData []struct {
	Strike    string `json:"strike"`
	Symbol    string `json:"symbol"`
	Direction string `json:"direction"`
}

func (data InstrumentData) GetSymbol(direction string) (symbol string, err error) {
	for _, v := range data {
		if v.Strike == "SPT" && v.Direction == direction {
			return v.Symbol, nil
		}
	}

	return "", fmt.Errorf("no symbol found for direction %s", direction)
}
