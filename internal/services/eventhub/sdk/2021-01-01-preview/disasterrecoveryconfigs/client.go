package disasterrecoveryconfigs

import "github.com/Azure/go-autorest/autorest"

type DisasterRecoveryConfigsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDisasterRecoveryConfigsClientWithBaseURI(endpoint string) DisasterRecoveryConfigsClient {
	return DisasterRecoveryConfigsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
