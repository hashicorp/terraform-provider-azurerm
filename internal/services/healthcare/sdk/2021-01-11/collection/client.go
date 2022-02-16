package collection

import "github.com/Azure/go-autorest/autorest"

type CollectionClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCollectionClientWithBaseURI(endpoint string) CollectionClient {
	return CollectionClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
