package clusters

import "github.com/Azure/go-autorest/autorest"

type ClustersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewClustersClientWithBaseURI(endpoint string) ClustersClient {
	return ClustersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
