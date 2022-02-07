package powerbidedicated

import "github.com/Azure/go-autorest/autorest"

type PowerBIDedicatedClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPowerBIDedicatedClientWithBaseURI(endpoint string) PowerBIDedicatedClient {
	return PowerBIDedicatedClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
