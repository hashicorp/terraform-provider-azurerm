package checknameavailability

import "github.com/Azure/go-autorest/autorest"

type CheckNameAvailabilityClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCheckNameAvailabilityClientWithBaseURI(endpoint string) CheckNameAvailabilityClient {
	return CheckNameAvailabilityClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
