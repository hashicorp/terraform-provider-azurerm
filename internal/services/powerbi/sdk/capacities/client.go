package capacities

import "github.com/Azure/go-autorest/autorest"

type CapacitiesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCapacitiesClientWithBaseURI(endpoint string) CapacitiesClient {
	return CapacitiesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
