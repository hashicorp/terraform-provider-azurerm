package origingroups

import "github.com/Azure/go-autorest/autorest"

type OriginGroupsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewOriginGroupsClientWithBaseURI(endpoint string) OriginGroupsClient {
	return OriginGroupsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
