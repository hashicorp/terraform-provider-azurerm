package hybridconnections

import "github.com/Azure/go-autorest/autorest"

type HybridConnectionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewHybridConnectionsClientWithBaseURI(endpoint string) HybridConnectionsClient {
	return HybridConnectionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
