package broker

import (
	"encoding/json"
)

type LoginData struct {
	Email    string  `json:"identifier"`
	Password string  `json:"password"`
	Token    *string `json:"token,omitempty"`
}

func httpLogin(url string, session *Session) error {
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

func httpLogout(url string, session *Session) error {
	resp, err := session.PostFromStruct(url, nil, nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
