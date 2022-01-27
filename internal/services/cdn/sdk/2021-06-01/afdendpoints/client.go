package afdendpoints

import "github.com/Azure/go-autorest/autorest"

type AFDEndpointsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAFDEndpointsClientWithBaseURI(endpoint string) AFDEndpointsClient {
	return AFDEndpointsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
