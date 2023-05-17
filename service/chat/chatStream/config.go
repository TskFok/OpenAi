package chatStream

import (
	"github.com/TskFok/OpenAi/app/global"
	"net/http"
)

const (
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
		BaseURL:    global.VicunaUrl,
		OrgID:      "",
		authToken:  authToken,

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}
