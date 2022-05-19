package account

import "github.com/Azure/go-autorest/autorest"

type AccountClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAccountClientWithBaseURI(endpoint string) AccountClient {
	return AccountClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
