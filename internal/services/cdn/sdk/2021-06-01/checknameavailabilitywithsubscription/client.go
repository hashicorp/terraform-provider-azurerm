package checknameavailabilitywithsubscription

import "github.com/Azure/go-autorest/autorest"

type CheckNameAvailabilityWithSubscriptionClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCheckNameAvailabilityWithSubscriptionClientWithBaseURI(endpoint string) CheckNameAvailabilityWithSubscriptionClient {
	return CheckNameAvailabilityWithSubscriptionClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
