package databases

import "github.com/Azure/go-autorest/autorest"

type DatabasesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDatabasesClientWithBaseURI(endpoint string) DatabasesClient {
	return DatabasesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
