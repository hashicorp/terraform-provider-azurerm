package managedidentity

import "github.com/Azure/go-autorest/autorest"

type ManagedIdentityClient struct {
	Client  autorest.Client
	baseUri string
}

func NewManagedIdentityClientWithBaseURI(endpoint string) ManagedIdentityClient {
	return ManagedIdentityClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
