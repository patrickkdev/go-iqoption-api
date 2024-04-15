package brokerws

import (
	"time"

	"github.com/patrickkdev/Go-IQOption-API/btypes"
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

type topAssetsResponse struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`
	Msg       struct {
		InstrumentType string        `json:"instrument_type"`
		RegionID       int           `json:"region_id"`
		Data           btypes.Assets `json:"data"`
	} `json:"msg"`
	Status int `json:"status"`
}

func (ws *Socket) GetTopAssets(type_ btypes.AssetType, timeout time.Time) (btypes.Assets, error) {
	requestEvent := &btypes.RequestEvent{
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

	resp, err := ws.EmitWithResponse(requestEvent, "top-assets", timeout)
	if err != nil {
		return btypes.Assets{}, err
	}

	responseEvent, err := tjson.Unmarshal[topAssetsResponse](resp)
	if err != nil {
		return btypes.Assets{}, err
	}

	return responseEvent.Msg.Data, nil
}
