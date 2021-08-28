package eventhubsclusters

import "github.com/Azure/go-autorest/autorest"

type EventHubsClustersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewEventHubsClustersClientWithBaseURI(endpoint string) EventHubsClustersClient {
	return EventHubsClustersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
