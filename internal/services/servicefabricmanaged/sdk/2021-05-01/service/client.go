package service

import "github.com/Azure/go-autorest/autorest"

type ServiceClient struct {
	Client  autorest.Client
	baseUri string
}

func NewServiceClientWithBaseURI(endpoint string) ServiceClient {
	return ServiceClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
