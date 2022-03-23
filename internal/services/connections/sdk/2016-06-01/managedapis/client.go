package managedapis

import "github.com/Azure/go-autorest/autorest"

type ManagedAPIsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewManagedAPIsClientWithBaseURI(endpoint string) ManagedAPIsClient {
	return ManagedAPIsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
