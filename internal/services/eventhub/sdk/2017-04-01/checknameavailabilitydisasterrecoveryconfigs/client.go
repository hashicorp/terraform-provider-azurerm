package checknameavailabilitydisasterrecoveryconfigs

import "github.com/Azure/go-autorest/autorest"

type CheckNameAvailabilityDisasterRecoveryConfigsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCheckNameAvailabilityDisasterRecoveryConfigsClientWithBaseURI(endpoint string) CheckNameAvailabilityDisasterRecoveryConfigsClient {
	return CheckNameAvailabilityDisasterRecoveryConfigsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
