package privateclouds

import "github.com/Azure/go-autorest/autorest"

type PrivateCloudsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPrivateCloudsClientWithBaseURI(endpoint string) PrivateCloudsClient {
	return PrivateCloudsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
