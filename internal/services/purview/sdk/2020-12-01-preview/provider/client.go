package provider

import "github.com/Azure/go-autorest/autorest"

type ProviderClient struct {
	Client  autorest.Client
	baseUri string
}

func NewProviderClientWithBaseURI(endpoint string) ProviderClient {
	return ProviderClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
