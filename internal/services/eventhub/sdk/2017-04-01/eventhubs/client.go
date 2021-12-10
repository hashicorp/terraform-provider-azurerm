package eventhubs

import "github.com/Azure/go-autorest/autorest"

type EventHubsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewEventHubsClientWithBaseURI(endpoint string) EventHubsClient {
	return EventHubsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
