package attestationproviders

import "github.com/Azure/go-autorest/autorest"

type AttestationProvidersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAttestationProvidersClientWithBaseURI(endpoint string) AttestationProvidersClient {
	return AttestationProvidersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
