package partnerregistrations

import "github.com/Azure/go-autorest/autorest"

type PartnerRegistrationsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPartnerRegistrationsClientWithBaseURI(endpoint string) PartnerRegistrationsClient {
	return PartnerRegistrationsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
