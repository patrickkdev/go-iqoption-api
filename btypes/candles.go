package btypes

type Candle struct {
	At     int64   `json:"at"`
	Close  float64 `json:"close"`
	From   int     `json:"from"`
	ID     int     `json:"id"`
	Max    float64 `json:"max"`
	Min    float64 `json:"min"`
	Open   float64 `json:"open"`
	To     int     `json:"to"`
	Volume int     `json:"volume"`
}

type Candles []Candle

func (candles Candles) GetLast() (candle Candle) {
	return candles[len(candles)-1]
}
