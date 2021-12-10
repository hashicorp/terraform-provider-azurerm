package namespaces

import "github.com/Azure/go-autorest/autorest"

type NamespacesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewNamespacesClientWithBaseURI(endpoint string) NamespacesClient {
	return NamespacesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
