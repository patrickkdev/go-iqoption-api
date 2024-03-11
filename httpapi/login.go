package httpapi

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

	return nil
}