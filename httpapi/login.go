package httpapi

import (
	"encoding/json"
	"patrickkdev/Go-IQOption-API/debug"
)

type LoginData struct {
	Identifier string  `json:"identifier"`
	Password   string  `json:"password"`
	Token      *string `json:"token,omitempty"`
}

func Login(url string, session *Session, data *LoginData) error {
	resp, err := session.PostFromStruct(url, data, nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&session)

	if err != nil {
		return err
	}

	debug.IfVerbose.PrintAsJSON(session)

	return nil
}