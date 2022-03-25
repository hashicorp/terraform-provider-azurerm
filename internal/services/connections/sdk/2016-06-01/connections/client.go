package connections

import "github.com/Azure/go-autorest/autorest"

type ConnectionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewConnectionsClientWithBaseURI(endpoint string) ConnectionsClient {
	return ConnectionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
