package servers

import "github.com/Azure/go-autorest/autorest"

type ServersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewServersClientWithBaseURI(endpoint string) ServersClient {
	return ServersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
