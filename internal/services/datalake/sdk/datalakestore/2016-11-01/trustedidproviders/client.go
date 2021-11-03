package trustedidproviders

import "github.com/Azure/go-autorest/autorest"

type TrustedIdProvidersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewTrustedIdProvidersClientWithBaseURI(endpoint string) TrustedIdProvidersClient {
	return TrustedIdProvidersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
