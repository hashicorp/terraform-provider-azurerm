package privatelinkresource

import "github.com/Azure/go-autorest/autorest"

type PrivateLinkResourceClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPrivateLinkResourceClientWithBaseURI(endpoint string) PrivateLinkResourceClient {
	return PrivateLinkResourceClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
