package configurations

import "github.com/Azure/go-autorest/autorest"

type ConfigurationsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewConfigurationsClientWithBaseURI(endpoint string) ConfigurationsClient {
	return ConfigurationsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
