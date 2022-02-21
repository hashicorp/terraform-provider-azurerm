package resource

import "github.com/Azure/go-autorest/autorest"

type ResourceClient struct {
	Client  autorest.Client
	baseUri string
}

func NewResourceClientWithBaseURI(endpoint string) ResourceClient {
	return ResourceClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
