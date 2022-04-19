package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client represents the data API client to use for feeding
type Client struct {
	address string
	key     string
	client  *http.Client
}

// New creates a new data API client
func New(address, key string) *Client {
	for strings.HasSuffix(address, "/") {
		address = strings.TrimSuffix(address, "/")
	}
	return &Client{
		address: address,
		key:     key,
		client:  &http.Client{},
	}
}

func (client *Client) execute(request *http.Request) (*http.Response, []byte, error) {
	request.Header.Add("Authorization", "Bearer "+client.key)
	response, err := client.client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		var errResponse *APIErrorResponse
		if err := json.Unmarshal(body, &errResponse); err == nil {
			return nil, nil, errResponse
		}
		return nil, nil, fmt.Errorf("HTTP status %d: %s", response.StatusCode, string(body))
	}
	return response, body, nil
}
