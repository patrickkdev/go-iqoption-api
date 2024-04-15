package brokerhttp

func Logout(url string, session *Session) error {
	resp, err := session.PostFromStruct(url, nil, nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
