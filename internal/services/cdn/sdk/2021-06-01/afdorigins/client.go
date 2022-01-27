package afdorigins

import "github.com/Azure/go-autorest/autorest"

type AFDOriginsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAFDOriginsClientWithBaseURI(endpoint string) AFDOriginsClient {
	return AFDOriginsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
