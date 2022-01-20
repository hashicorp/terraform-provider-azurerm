package frontdoors

import "github.com/Azure/go-autorest/autorest"

type FrontDoorsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewFrontDoorsClientWithBaseURI(endpoint string) FrontDoorsClient {
	return FrontDoorsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
