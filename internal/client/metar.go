package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// FeedMETARs feeds METARs into the data server
func (client *Client) FeedMETARs(metars []string) ([]int, error) {
	data, err := json.Marshal(map[string][]string{
		"data": metars,
	})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, client.address+endpointMETARs, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	_, body, err := client.execute(request)
	if err != nil {
		return nil, err
	}

	responseData := new(struct {
		Duplicates []int `json:"duplicates"`
	})
	if err := json.Unmarshal(body, responseData); err != nil {
		return nil, err
	}

	return responseData.Duplicates, nil
}
