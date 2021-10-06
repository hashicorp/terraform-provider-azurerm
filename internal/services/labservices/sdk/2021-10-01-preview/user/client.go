package user

import "github.com/Azure/go-autorest/autorest"

type UserClient struct {
	Client  autorest.Client
	baseUri string
}

func NewUserClientWithBaseURI(endpoint string) UserClient {
	return UserClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
