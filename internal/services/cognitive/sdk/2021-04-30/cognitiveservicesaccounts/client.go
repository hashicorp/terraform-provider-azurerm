package cognitiveservicesaccounts

import "github.com/Azure/go-autorest/autorest"

type CognitiveServicesAccountsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCognitiveServicesAccountsClientWithBaseURI(endpoint string) CognitiveServicesAccountsClient {
	return CognitiveServicesAccountsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
