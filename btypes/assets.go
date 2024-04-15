package btypes

type AssetType string

const (
	AssetTypeTurbo   AssetType = "turbo-option"
	AssetTypeBinary  AssetType = "binary-option"
	AssetTypeDigital AssetType = "digital-option"
)

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

func (assets *Assets) FilterOpen() Assets {
	var newAssets Assets = Assets{}

	for _, asset := range *assets {
		if asset.Expiration > 0 &&
			asset.Volume > 0 &&
			asset.SpotProfit > 0 &&
			asset.Volatility > 0 {
			newAssets = append(newAssets, asset)
		}
	}

	return newAssets
}

func (assets *Assets) RemoveByID(id int) {
	for i, asset := range *assets {
		if asset.ActiveID == id {
			(*assets)[i] = (*assets)[len(*assets)-1]
			*assets = (*assets)[:len(*assets)-1]
			break
		}
	}
}
