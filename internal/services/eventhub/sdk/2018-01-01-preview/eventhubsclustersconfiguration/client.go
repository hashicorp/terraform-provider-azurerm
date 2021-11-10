package eventhubsclustersconfiguration

import "github.com/Azure/go-autorest/autorest"

type EventHubsClustersConfigurationClient struct {
	Client  autorest.Client
	baseUri string
}

func NewEventHubsClustersConfigurationClientWithBaseURI(endpoint string) EventHubsClustersConfigurationClient {
	return EventHubsClustersConfigurationClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
