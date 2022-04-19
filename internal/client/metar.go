package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
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

// FeedMETARsRelaxed works the same as FeedMETARs with the exception that it simply skips METARs with an invalid format
func (client *Client) FeedMETARsRelaxed(metars []string) error {
	_, err := client.FeedMETARs(metars)
	if err != nil {
		var errResponse *APIErrorResponse
		if !errors.As(err, &errResponse) || len(errResponse.Errors) == 0 {
			return err
		}
		if errResponse.Errors[0].Type != "data.metars.invalidFormat" {
			return err
		}
		asshole := int(errResponse.Errors[0].Details["index"].(float64))
		log.Warn().Msgf("skipping invalid METAR: '%s'", metars[asshole])
		metars[asshole] = metars[len(metars)-1]
		return client.FeedMETARsRelaxed(metars[:len(metars)-1])
	}
	return nil
}
