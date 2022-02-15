package locations

import "github.com/Azure/go-autorest/autorest"

type LocationsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewLocationsClientWithBaseURI(endpoint string) LocationsClient {
	return LocationsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
