package secrets

import "github.com/Azure/go-autorest/autorest"

type SecretsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSecretsClientWithBaseURI(endpoint string) SecretsClient {
	return SecretsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
