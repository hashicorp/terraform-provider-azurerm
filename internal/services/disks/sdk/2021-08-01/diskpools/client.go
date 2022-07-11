package diskpools

import "github.com/Azure/go-autorest/autorest"

type DiskPoolsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDiskPoolsClientWithBaseURI(endpoint string) DiskPoolsClient {
	return DiskPoolsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
