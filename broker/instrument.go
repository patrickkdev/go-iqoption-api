package broker

import (
	"fmt"

	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
	"github.com/patrickkdev/Go-IQOption-API/internal/types"
)

func (c *Client) getInstrument(activeID int, exp int) (*types.Instrument, error) {
	requestEvent := requestEvent{
		Name: "sendMessage",
		Msg: map[string]interface{}{
			"name":    "digital-option-instruments.get-instruments",
			"version": "2.0",
			"body": map[string]interface{}{
				"instrument_type": "digital-option",
				"asset_id":        activeID,
			},
		},
	}.WithRandomRequestId()

	resp, err := c.ws.EmitWithResponse(requestEvent, "instruments", c.getTimeout())
	if err != nil {
		return nil, err
	}

	responseEvent, err := tjson.Unmarshal[types.InstrumentsResponseEvent](resp)
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

		debug.IfVerbose.Printf("Instrument expiration: %d\n", instrument.Expiration)
	}

	return nil, fmt.Errorf("no instrument found for active id %d and expiration %d", activeID, exp)
}
