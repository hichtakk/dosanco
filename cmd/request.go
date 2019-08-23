package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type responseMessage struct {
	Message string `json:"message"`
}

func sendRequest(method string, url string, reqJSON []byte) ([]byte, error) {
	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer(reqJSON),
	)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	if resp.StatusCode != 200 {
		return body, fmt.Errorf(string(resp.StatusCode))
	}

	return body, nil
}
