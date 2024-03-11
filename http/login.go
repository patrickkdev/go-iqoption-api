package httpapi

import (
	httpapi "github.com/patrickkdev/Go-IQOption-API/http/utils"
)

type LoginData struct {
	Identifier string  `json:"identifier"`
	Password   string  `json:"password"`
	Token      *string `json:"token,omitempty"`
}

func Login(url string, data LoginData) error {
	resp, err := httpapi.PostJson(url, data)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
