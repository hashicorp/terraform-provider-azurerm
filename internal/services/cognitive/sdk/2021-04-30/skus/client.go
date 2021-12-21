package skus

import "github.com/Azure/go-autorest/autorest"

type SkusClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSkusClientWithBaseURI(endpoint string) SkusClient {
	return SkusClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
