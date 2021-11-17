package query

import "github.com/Azure/go-autorest/autorest"

type QueryClient struct {
	Client  autorest.Client
	baseUri string
}

func NewQueryClientWithBaseURI(endpoint string) QueryClient {
	return QueryClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
