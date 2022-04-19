package client

import (
	"encoding/json"
	"net/http"
)

// KeyInfo represents the important information of a data API key
type KeyInfo struct {
	Quota        int64 `json:"quota"`
	RateLimit    int   `json:"rate_limit"`
	Capabilities uint  `json:"capabilities"`
}

// GetKeyInfo retrieves the KeyInfo about the currently used API key
func (client *Client) GetKeyInfo() (*KeyInfo, error) {
	request, err := http.NewRequest(http.MethodGet, client.address+endpointKeyInfo, nil)
	if err != nil {
		return nil, err
	}

	_, body, err := client.execute(request)
	if err != nil {
		return nil, err
	}

	info := new(KeyInfo)
	if err := json.Unmarshal(body, info); err != nil {
		return nil, err
	}

	return info, nil
}
