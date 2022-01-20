package tenants

import "github.com/Azure/go-autorest/autorest"

type TenantsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewTenantsClientWithBaseURI(endpoint string) TenantsClient {
	return TenantsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
