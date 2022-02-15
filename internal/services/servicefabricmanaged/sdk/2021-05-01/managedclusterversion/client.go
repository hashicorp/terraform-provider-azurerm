package managedclusterversion

import "github.com/Azure/go-autorest/autorest"

type ManagedClusterVersionClient struct {
	Client  autorest.Client
	baseUri string
}

func NewManagedClusterVersionClientWithBaseURI(endpoint string) ManagedClusterVersionClient {
	return ManagedClusterVersionClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
