package broker

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Session struct {
	Code            string    `json:"code"`
	SSID            string    `json:"ssid"`
	ClientSessionID string    `json:"client_session_id"`
	LoginData       LoginData `json:"login_data,omitempty"`
}

func NewSession(loginData LoginData) Session {
	return Session{LoginData: loginData}
}

func (sD *Session) PostFromStruct(url string, data interface{}, customHeaders map[string]string) (*http.Response, error) {
	var jsonBody []byte = []byte{}
	if data != nil {
		var err error // Declare error for later use
		jsonBody, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	for key, value := range customHeaders {
		request.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
