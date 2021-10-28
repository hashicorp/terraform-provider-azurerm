package managedhsms

import "github.com/Azure/go-autorest/autorest"

type ManagedHsmsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewManagedHsmsClientWithBaseURI(endpoint string) ManagedHsmsClient {
	return ManagedHsmsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
