package rules

import "github.com/Azure/go-autorest/autorest"

type RulesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRulesClientWithBaseURI(endpoint string) RulesClient {
	return RulesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
