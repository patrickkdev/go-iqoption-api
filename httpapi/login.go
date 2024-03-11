package httpapi

import utils "github.com/patrickkdev/Go-IQOption-API/httpapi/utils"

type LoginData struct {
	Identifier string  `json:"identifier"`
	Password   string  `json:"password"`
	Token      *string `json:"token,omitempty"`
}

func Login(url string, data *LoginData) error {
	resp, err := utils.PostFromStruct(url, data)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
