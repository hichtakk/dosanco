package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hichtakk/dosanco/handler"
)

type responseMessage struct {
	Message string `json:"message"`
}

func sendRequest(method string, url string, reqJSON []byte) ([]byte, error) {
	// prepare
	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer(reqJSON),
	)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// do request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if debug != false {
		fmt.Printf("[%d] %s %s\n", resp.StatusCode, resp.Request.Method, resp.Request.URL)
	}

	// read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("read response error")
	}
	if resp.StatusCode != 200 {
		errBody := new(handler.ErrorResponse)
		if err := json.Unmarshal(body, errBody); err != nil {
			return []byte{}, fmt.Errorf("unmarshall response error")
		}
		msg := ""
		if errBody.Message != "" {
			msg = errBody.Message
		} else {
			msg = errBody.Error.Message
		}
		return []byte{}, fmt.Errorf(msg)
	}

	return body, nil
}
