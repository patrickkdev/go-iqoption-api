package broker

import (
	"fmt"
	"slices"

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

	if responseEvent.Status != 2000 {
		return Assets{}, fmt.Errorf("error getting assets")
	}

	return responseEvent.Msg.Data, nil
}

// Returns a copy of the assets without the asset with the given id
func (assets Assets) RemoveByID(id int) Assets {
	index := slices.IndexFunc(assets, func(asset types.Asset) bool {
		return asset.ActiveID == id
	})

	if index == -1 {
		return assets
	}

	return assets.RemoveByIndex(index)
}

func (assets Assets) RemoveByIndex(index int) Assets {
	if index < 0 || index >= len(assets) {
		fmt.Printf("RemoveByIndex: Invalid index: %d\n", index)

		return assets
	}

	assets[index] = assets[len(assets)-1]
	return assets[:len(assets)-1]
}

// Less performant
//
// func filter(assets Assets, keepCondition func(types.Asset) bool) Assets {
// 	newAssets := Assets{}

// 	for _, asset := range assets {
// 		if keepCondition(asset) {
// 			newAssets = append(newAssets, asset)
// 		}
// 	}

// 	return newAssets
// }
