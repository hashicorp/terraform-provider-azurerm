package origins

import "github.com/Azure/go-autorest/autorest"

type OriginsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewOriginsClientWithBaseURI(endpoint string) OriginsClient {
	return OriginsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
