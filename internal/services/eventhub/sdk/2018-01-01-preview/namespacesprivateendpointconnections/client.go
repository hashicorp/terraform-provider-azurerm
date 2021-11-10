package namespacesprivateendpointconnections

import "github.com/Azure/go-autorest/autorest"

type NamespacesPrivateEndpointConnectionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewNamespacesPrivateEndpointConnectionsClientWithBaseURI(endpoint string) NamespacesPrivateEndpointConnectionsClient {
	return NamespacesPrivateEndpointConnectionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
