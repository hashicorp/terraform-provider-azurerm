package afdcustomdomains

import "github.com/Azure/go-autorest/autorest"

type AFDCustomDomainsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAFDCustomDomainsClientWithBaseURI(endpoint string) AFDCustomDomainsClient {
	return AFDCustomDomainsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
