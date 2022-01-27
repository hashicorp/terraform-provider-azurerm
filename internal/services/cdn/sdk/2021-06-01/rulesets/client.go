package rulesets

import "github.com/Azure/go-autorest/autorest"

type RuleSetsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRuleSetsClientWithBaseURI(endpoint string) RuleSetsClient {
	return RuleSetsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
