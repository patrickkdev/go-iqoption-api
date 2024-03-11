package httpapi

import utils "github.com/patrickkdev/Go-IQOption-API/httpapi/utils"

func Logout(url string) error {
	resp, err := utils.PostFromStruct(url, nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
