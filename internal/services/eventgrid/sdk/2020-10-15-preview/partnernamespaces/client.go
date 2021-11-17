package partnernamespaces

import "github.com/Azure/go-autorest/autorest"

type PartnerNamespacesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPartnerNamespacesClientWithBaseURI(endpoint string) PartnerNamespacesClient {
	return PartnerNamespacesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
