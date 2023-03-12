package testhelper

import "github.com/sashabaranov/go-openai"

const (
	testAuthToken = ""
)

// NewTestClient :
func NewTestClient(url string) *openai.Client {
	conf := openai.DefaultConfig(testAuthToken)
	conf.BaseURL = url

	return openai.NewClientWithConfig(conf)
}
