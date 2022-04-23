package privatelinkresources

import "github.com/Azure/go-autorest/autorest"

type PrivateLinkResourcesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPrivateLinkResourcesClientWithBaseURI(endpoint string) PrivateLinkResourcesClient {
	return PrivateLinkResourcesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
