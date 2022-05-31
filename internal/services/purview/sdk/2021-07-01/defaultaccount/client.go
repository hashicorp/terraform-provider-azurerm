package defaultaccount

import "github.com/Azure/go-autorest/autorest"

type DefaultAccountClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDefaultAccountClientWithBaseURI(endpoint string) DefaultAccountClient {
	return DefaultAccountClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
