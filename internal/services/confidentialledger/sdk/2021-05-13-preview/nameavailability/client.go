package nameavailability

import "github.com/Azure/go-autorest/autorest"

type NameAvailabilityClient struct {
	Client  autorest.Client
	baseUri string
}

func NewNameAvailabilityClientWithBaseURI(endpoint string) NameAvailabilityClient {
	return NameAvailabilityClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
