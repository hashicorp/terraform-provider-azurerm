package endpoints

import "github.com/Azure/go-autorest/autorest"

type EndpointsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewEndpointsClientWithBaseURI(endpoint string) EndpointsClient {
	return EndpointsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
