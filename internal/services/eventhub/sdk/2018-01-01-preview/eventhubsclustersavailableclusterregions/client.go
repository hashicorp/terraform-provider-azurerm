package eventhubsclustersavailableclusterregions

import "github.com/Azure/go-autorest/autorest"

type EventHubsClustersAvailableClusterRegionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewEventHubsClustersAvailableClusterRegionsClientWithBaseURI(endpoint string) EventHubsClustersAvailableClusterRegionsClient {
	return EventHubsClustersAvailableClusterRegionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
