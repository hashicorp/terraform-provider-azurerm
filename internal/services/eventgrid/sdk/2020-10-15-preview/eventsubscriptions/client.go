package eventsubscriptions

import "github.com/Azure/go-autorest/autorest"

type EventSubscriptionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewEventSubscriptionsClientWithBaseURI(endpoint string) EventSubscriptionsClient {
	return EventSubscriptionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
