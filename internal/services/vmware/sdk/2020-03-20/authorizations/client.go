package authorizations

import "github.com/Azure/go-autorest/autorest"

type AuthorizationsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAuthorizationsClientWithBaseURI(endpoint string) AuthorizationsClient {
	return AuthorizationsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
