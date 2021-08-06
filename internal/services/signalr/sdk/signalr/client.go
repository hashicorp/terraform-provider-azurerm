package signalr

import "github.com/Azure/go-autorest/autorest"

type SignalRClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSignalRClientWithBaseURI(endpoint string) SignalRClient {
	return SignalRClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
