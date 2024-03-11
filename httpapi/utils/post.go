package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func PostFromStruct(url string, data interface{}) (*http.Response, error) {
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

	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
