package mhsmlistprivateendpointconnections

import "github.com/Azure/go-autorest/autorest"

type MHSMListPrivateEndpointConnectionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewMHSMListPrivateEndpointConnectionsClientWithBaseURI(endpoint string) MHSMListPrivateEndpointConnectionsClient {
	return MHSMListPrivateEndpointConnectionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
