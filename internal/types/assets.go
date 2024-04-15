package types

type TopAssetsResponse struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`
	Msg       struct {
		InstrumentType string  `json:"instrument_type"`
		RegionID       int     `json:"region_id"`
		Data           []Asset `json:"data"`
	} `json:"msg"`
	Status int `json:"status"`
}

type Asset struct {
	ActiveID    int     `json:"active_id"`
	Spread      float64 `json:"spread,omitempty"`
	Diff5Min    float64 `json:"diff5_min,omitempty"`
	DiffHour    float64 `json:"diff_hour,omitempty"`
	DiffMonth   float64 `json:"diff_month,omitempty"`
	CurPrice    float64 `json:"cur_price,omitempty"`
	Volume      float64 `json:"volume,omitempty"`
	Popularity  float64 `json:"popularity,omitempty"`
	Expiration  int64   `json:"expiration"`
	SpotProfit  float64 `json:"spot_profit,omitempty"`
	TradersMood float64 `json:"traders_mood"`
	DiffMin     float64 `json:"diff_min,omitempty"`
	DiffDay     float64 `json:"diff_day,omitempty"`
	Volatility  float64 `json:"volatility,omitempty"`
}
