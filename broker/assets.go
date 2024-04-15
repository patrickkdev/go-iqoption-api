package broker

import (
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
	"github.com/patrickkdev/Go-IQOption-API/internal/types"
)

type AssetType string

const (
	AssetTypeTurbo   AssetType = "turbo-option"
	AssetTypeBinary  AssetType = "binary-option"
	AssetTypeDigital AssetType = "digital-option"
)

type Assets []types.Asset

// Gets available assets (or 'pairs' like 'EUR/USD') for a given asset type like 'binary-option' or 'digital-option'
// Pairs returned are not garanteed to be tradable.
// Calling .FilterOpen() on the returned assets will TRY to filter out non-tradable pairs
func (c *Client) GetAssets(type_ AssetType) (Assets, error) {
	requestEvent := requestEvent{
		Name: "sendMessage",
		Msg: map[string]interface{}{
			"name":    "get-top-assets",
			"version": "3.0",
			"body": map[string]interface{}{
				"instrument_type": string(type_),
				"region_id":       -1,
			},
		},
	}.WithRandomRequestId()

	resp, err := c.ws.EmitWithResponse(requestEvent, "top-assets", c.getTimeout())
	if err != nil {
		return Assets{}, err
	}

	responseEvent, err := tjson.Unmarshal[types.TopAssetsResponse](resp)
	if err != nil {
		return Assets{}, err
	}

	return responseEvent.Msg.Data, nil
}

// Filters out assets that are probably not tradable
func (assets *Assets) WithoutNonTradable() Assets {
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

// Returns a copy of the assets without the asset with the given id
func (assets *Assets) Without(id int) Assets {
	newAssets := Assets{}

	for _, asset := range *assets {
		if asset.ActiveID != id {
			newAssets = append(newAssets, asset)
		}
	}

	return newAssets
}
