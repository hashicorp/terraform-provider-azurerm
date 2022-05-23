package servicelinker

import "github.com/Azure/go-autorest/autorest"

type ServiceLinkerClient struct {
	Client  autorest.Client
	baseUri string
}

func NewServiceLinkerClientWithBaseURI(endpoint string, id string) ServiceLinkerClient {
	return ServiceLinkerClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
