package resourceskus

import "github.com/Azure/go-autorest/autorest"

type ResourceSkusClient struct {
	Client  autorest.Client
	baseUri string
}

func NewResourceSkusClientWithBaseURI(endpoint string) ResourceSkusClient {
	return ResourceSkusClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
