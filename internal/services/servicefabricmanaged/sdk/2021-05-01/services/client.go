package services

import "github.com/Azure/go-autorest/autorest"

type ServicesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewServicesClientWithBaseURI(endpoint string) ServicesClient {
	return ServicesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
