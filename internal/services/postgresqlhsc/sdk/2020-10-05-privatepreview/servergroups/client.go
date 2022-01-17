package servergroups

import "github.com/Azure/go-autorest/autorest"

type ServerGroupsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewServerGroupsClientWithBaseURI(endpoint string) ServerGroupsClient {
	return ServerGroupsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
