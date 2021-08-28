package networkrulesets

import "github.com/Azure/go-autorest/autorest"

type NetworkRuleSetsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewNetworkRuleSetsClientWithBaseURI(endpoint string) NetworkRuleSetsClient {
	return NetworkRuleSetsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
