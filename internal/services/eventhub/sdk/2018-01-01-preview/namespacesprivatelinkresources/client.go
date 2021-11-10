package namespacesprivatelinkresources

import "github.com/Azure/go-autorest/autorest"

type NamespacesPrivateLinkResourcesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewNamespacesPrivateLinkResourcesClientWithBaseURI(endpoint string) NamespacesPrivateLinkResourcesClient {
	return NamespacesPrivateLinkResourcesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
