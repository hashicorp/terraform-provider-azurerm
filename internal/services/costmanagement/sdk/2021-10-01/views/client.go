package views

import "github.com/Azure/go-autorest/autorest"

type ViewsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewViewsClientWithBaseURI(endpoint string) ViewsClient {
	return ViewsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
