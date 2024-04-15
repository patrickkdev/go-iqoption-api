package brokerhttp

import (
	"encoding/json"
)

type LoginData struct {
	Email    string  `json:"identifier"`
	Password string  `json:"password"`
	Token    *string `json:"token,omitempty"`
}

func Login(url string, session *Session) error {
	resp, err := session.PostFromStruct(url, session.LoginData, nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(session)

	if err != nil {
		return err
	}

	return nil
}
