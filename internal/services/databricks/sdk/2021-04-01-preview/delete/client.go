package delete

import "github.com/Azure/go-autorest/autorest"

type DELETEClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDELETEClientWithBaseURI(endpoint string) DELETEClient {
	return DELETEClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
