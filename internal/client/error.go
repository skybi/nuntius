package client

import "strings"

// APIErrorResponse represents an error response following the error structure of the data server
type APIErrorResponse struct {
	Status int         `json:"status"`
	Errors []*APIError `json:"errors"`
}

func (resp *APIErrorResponse) Error() string {
	msgs := make([]string, len(resp.Errors))
	for i, err := range resp.Errors {
		msgs[i] = err.Message
	}
	return strings.Join(msgs, ", ")
}

// APIError represents a single error present in an APIErrorResponse
type APIError struct {
	Type    string         `json:"type"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}
