package schemaregistry

import "github.com/Azure/go-autorest/autorest"

type SchemaRegistryClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSchemaRegistryClientWithBaseURI(endpoint string) SchemaRegistryClient {
	return SchemaRegistryClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
