package roles

import "github.com/Azure/go-autorest/autorest"

type RolesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRolesClientWithBaseURI(endpoint string) RolesClient {
	return RolesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
