package wsapi

func GetUserProfileClient(ws *WSocket, userId int) error {
	event := map[string]interface{}{
		"name": "get-user-profile-client",
		"body": map[string]interface{}{
			"user_id": userId,
		},
		"version": "1.0",
	}

	ws.Write(event)

	return nil
}