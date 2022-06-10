package checknameavailabilitynamespaces

import "github.com/Azure/go-autorest/autorest"

type CheckNameAvailabilityNamespacesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCheckNameAvailabilityNamespacesClientWithBaseURI(endpoint string) CheckNameAvailabilityNamespacesClient {
	return CheckNameAvailabilityNamespacesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
