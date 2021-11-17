package eventchannels

import "github.com/Azure/go-autorest/autorest"

type EventChannelsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewEventChannelsClientWithBaseURI(endpoint string) EventChannelsClient {
	return EventChannelsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
