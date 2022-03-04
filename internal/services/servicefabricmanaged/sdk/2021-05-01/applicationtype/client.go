package applicationtype

import "github.com/Azure/go-autorest/autorest"

type ApplicationTypeClient struct {
	Client  autorest.Client
	baseUri string
}

func NewApplicationTypeClientWithBaseURI(endpoint string) ApplicationTypeClient {
	return ApplicationTypeClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
