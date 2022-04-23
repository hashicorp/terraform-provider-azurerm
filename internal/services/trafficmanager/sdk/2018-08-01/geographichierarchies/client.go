package geographichierarchies

import "github.com/Azure/go-autorest/autorest"

type GeographicHierarchiesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewGeographicHierarchiesClientWithBaseURI(endpoint string) GeographicHierarchiesClient {
	return GeographicHierarchiesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
