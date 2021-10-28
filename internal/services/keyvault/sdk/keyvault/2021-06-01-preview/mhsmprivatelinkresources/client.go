package mhsmprivatelinkresources

import "github.com/Azure/go-autorest/autorest"

type MHSMPrivateLinkResourcesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewMHSMPrivateLinkResourcesClientWithBaseURI(endpoint string) MHSMPrivateLinkResourcesClient {
	return MHSMPrivateLinkResourcesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
