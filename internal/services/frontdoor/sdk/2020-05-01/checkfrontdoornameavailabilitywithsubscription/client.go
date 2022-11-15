package checkfrontdoornameavailabilitywithsubscription

import "github.com/Azure/go-autorest/autorest"

type CheckFrontDoorNameAvailabilityWithSubscriptionClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCheckFrontDoorNameAvailabilityWithSubscriptionClientWithBaseURI(endpoint string) CheckFrontDoorNameAvailabilityWithSubscriptionClient {
	return CheckFrontDoorNameAvailabilityWithSubscriptionClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
