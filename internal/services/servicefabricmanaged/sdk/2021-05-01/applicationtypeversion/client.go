package applicationtypeversion

import "github.com/Azure/go-autorest/autorest"

type ApplicationTypeVersionClient struct {
	Client  autorest.Client
	baseUri string
}

func NewApplicationTypeVersionClientWithBaseURI(endpoint string) ApplicationTypeVersionClient {
	return ApplicationTypeVersionClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
