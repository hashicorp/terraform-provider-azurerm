package hcxenterprisesites

import "github.com/Azure/go-autorest/autorest"

type HcxEnterpriseSitesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewHcxEnterpriseSitesClientWithBaseURI(endpoint string) HcxEnterpriseSitesClient {
	return HcxEnterpriseSitesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
