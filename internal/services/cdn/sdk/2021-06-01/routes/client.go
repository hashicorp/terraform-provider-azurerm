package routes

import "github.com/Azure/go-autorest/autorest"

type RoutesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRoutesClientWithBaseURI(endpoint string) RoutesClient {
	return RoutesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
