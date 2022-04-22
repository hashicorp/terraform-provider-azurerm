package customapis

import "github.com/Azure/go-autorest/autorest"

type CustomAPIsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCustomAPIsClientWithBaseURI(endpoint string) CustomAPIsClient {
	return CustomAPIsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
