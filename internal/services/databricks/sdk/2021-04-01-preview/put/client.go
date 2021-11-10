package put

import "github.com/Azure/go-autorest/autorest"

type PUTClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPUTClientWithBaseURI(endpoint string) PUTClient {
	return PUTClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
