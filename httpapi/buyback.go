package httpapi

func Buyback(url string, session *Session, optionId int) error {
	data := struct {
		OptionId int `json:"option_id"`
	}{
		OptionId: optionId,
	}

	resp, err := session.PostFromStruct(url, data, nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}