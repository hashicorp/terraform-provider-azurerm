package managedcluster

import "github.com/Azure/go-autorest/autorest"

type ManagedClusterClient struct {
	Client  autorest.Client
	baseUri string
}

func NewManagedClusterClientWithBaseURI(endpoint string) ManagedClusterClient {
	return ManagedClusterClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
