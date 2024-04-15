package brokerws

import (
	"time"

	"github.com/patrickkdev/Go-IQOption-API/btypes"
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

func (ws *Socket) GetCoreProfile(timeout time.Time) (btypes.CoreProfile, error) {
	eventMsg := map[string]interface{}{
		"name":    "core.get-profile",
		"body":    map[string]interface{}{},
		"version": "1.0",
	}

	requestEvent := &btypes.RequestEvent{
		Name: "sendMessage",
		Msg:  eventMsg,
	}

	resp, err := ws.EmitWithResponse(requestEvent, "profile", timeout)
	if err != nil {
		return btypes.CoreProfile{}, err
	}

	responseEvent, err := tjson.Unmarshal[btypes.CoreProfile](resp)
	if err != nil {
		return btypes.CoreProfile{}, err
	}

	return responseEvent, nil
}

func (ws *Socket) GetUserProfileClient(userId int, timeout time.Time) (btypes.UserProfileClient, error) {
	eventMsg := map[string]interface{}{
		"name": "get-user-profile-client",
		"body": map[string]interface{}{
			"user_id": userId,
		},
		"version": "1.0",
	}

	event := &btypes.RequestEvent{
		Name: "sendMessage",
		Msg:  eventMsg,
	}

	resp, err := ws.EmitWithResponse(event, "user-profile-client", timeout)
	if err != nil {
		return btypes.UserProfileClient{}, err
	}

	responseEvent, err := tjson.Unmarshal[btypes.UserProfileClient](resp)
	if err != nil {
		return btypes.UserProfileClient{}, err
	}

	return responseEvent, nil
}
