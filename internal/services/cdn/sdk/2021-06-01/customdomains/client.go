package customdomains

import "github.com/Azure/go-autorest/autorest"

type CustomDomainsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCustomDomainsClientWithBaseURI(endpoint string) CustomDomainsClient {
	return CustomDomainsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
