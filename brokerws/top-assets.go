package brokerws

import (
	"time"

	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

type AssetType string

const (
	AssetTypeTurbo   AssetType = "turbo-option"
	AssetTypeBinary  AssetType = "binary-option"
	AssetTypeDigital AssetType = "digital-option"
)

type topAssetsResponse struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`
	Msg       struct {
		InstrumentType string `json:"instrument_type"`
		RegionID       int    `json:"region_id"`
		Data           Assets `json:"data"`
	} `json:"msg"`
	Status int `json:"status"`
}
type asset struct {
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

type Assets []asset

func (assets *Assets) FilterOpen() *Assets {
	var openAssets Assets
	for _, asset := range *assets {
		if 	asset.Expiration > 0 && 
				asset.Volume > 0 && 
				asset.SpotProfit > 0 &&
				asset.Volatility > 0 &&
				asset.ActiveID < 1000 {
			openAssets = append(openAssets, asset)
		}
	}

	return &openAssets
}

func GetTopAssets(ws *Socket, type_ AssetType, timeout time.Time) (*Assets, error) {
	requestEvent := &RequestEvent{
		Name: "sendMessage",
		Msg: map[string]interface{}{
			"name":    "get-top-assets",
			"version": "3.0",
			"body": map[string]interface{}{
				"instrument_type": string(type_),
				"region_id":       -1,
			},
		},
	}

	resp, err := EmitWithResponse(ws, requestEvent, "top-assets", timeout)
	if err != nil {
		return nil, err
	}

	responseEvent, err := tjson.Unmarshal[topAssetsResponse](resp)
	if err != nil {
		return nil, err
	}

	return &responseEvent.Msg.Data, nil
}
