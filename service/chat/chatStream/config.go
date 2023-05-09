package chatStream

import (
	"net/http"
)

const (
	apiURLv1                       = "http://region-9.seetacloud.com:41584/v1"
	defaultEmptyMessagesLimit uint = 300
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	authToken string

	HTTPClient *http.Client

	BaseURL string
	OrgID   string

	EmptyMessagesLimit uint
}

func DefaultConfig(authToken string) ClientConfig {
	return ClientConfig{
		HTTPClient: &http.Client{},
		BaseURL:    apiURLv1,
		OrgID:      "",
		authToken:  authToken,

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}
