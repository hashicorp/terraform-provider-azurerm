package links

import "github.com/Azure/go-autorest/autorest"

type LinksClient struct {
	Client  autorest.Client
	baseUri string
}

func NewLinksClientWithBaseURI(endpoint string) LinksClient {
	return LinksClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
