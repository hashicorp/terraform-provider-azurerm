package eventhubsclustersnamespace

import "github.com/Azure/go-autorest/autorest"

type EventHubsClustersNamespaceClient struct {
	Client  autorest.Client
	baseUri string
}

func NewEventHubsClustersNamespaceClientWithBaseURI(endpoint string) EventHubsClustersNamespaceClient {
	return EventHubsClustersNamespaceClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
