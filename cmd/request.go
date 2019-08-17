package cmd

import (
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
)

type responseMessage struct {
	Message string `json:"message"`
}

func sendRequest(method string, url string, reqJson []byte) ([]byte, error) {
	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer(reqJson),
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
