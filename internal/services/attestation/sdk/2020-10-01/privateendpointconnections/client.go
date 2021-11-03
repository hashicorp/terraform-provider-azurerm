package privateendpointconnections

import "github.com/Azure/go-autorest/autorest"

type PrivateEndpointConnectionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPrivateEndpointConnectionsClientWithBaseURI(endpoint string) PrivateEndpointConnectionsClient {
	return PrivateEndpointConnectionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
