package exports

import "github.com/Azure/go-autorest/autorest"

type ExportsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewExportsClientWithBaseURI(endpoint string) ExportsClient {
	return ExportsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
