package configurationstores

import "github.com/Azure/go-autorest/autorest"

type ConfigurationStoresClient struct {
	Client  autorest.Client
	baseUri string
}

func NewConfigurationStoresClientWithBaseURI(endpoint string) ConfigurationStoresClient {
	return ConfigurationStoresClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
