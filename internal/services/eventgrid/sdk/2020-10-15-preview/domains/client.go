package domains

import "github.com/Azure/go-autorest/autorest"

type DomainsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDomainsClientWithBaseURI(endpoint string) DomainsClient {
	return DomainsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
