package application

import "github.com/Azure/go-autorest/autorest"

type ApplicationClient struct {
	Client  autorest.Client
	baseUri string
}

func NewApplicationClientWithBaseURI(endpoint string) ApplicationClient {
	return ApplicationClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
