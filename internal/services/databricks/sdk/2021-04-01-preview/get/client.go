package get

import "github.com/Azure/go-autorest/autorest"

type GETClient struct {
	Client  autorest.Client
	baseUri string
}

func NewGETClientWithBaseURI(endpoint string) GETClient {
	return GETClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
