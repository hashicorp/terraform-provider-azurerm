package mhsmprivateendpointconnections

import "github.com/Azure/go-autorest/autorest"

type MHSMPrivateEndpointConnectionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewMHSMPrivateEndpointConnectionsClientWithBaseURI(endpoint string) MHSMPrivateEndpointConnectionsClient {
	return MHSMPrivateEndpointConnectionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
