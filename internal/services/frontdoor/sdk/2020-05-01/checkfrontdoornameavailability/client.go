package checkfrontdoornameavailability

import "github.com/Azure/go-autorest/autorest"

type CheckFrontDoorNameAvailabilityClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCheckFrontDoorNameAvailabilityClientWithBaseURI(endpoint string) CheckFrontDoorNameAvailabilityClient {
	return CheckFrontDoorNameAvailabilityClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
