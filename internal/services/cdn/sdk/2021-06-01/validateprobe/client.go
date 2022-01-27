package validateprobe

import "github.com/Azure/go-autorest/autorest"

type ValidateProbeClient struct {
	Client  autorest.Client
	baseUri string
}

func NewValidateProbeClientWithBaseURI(endpoint string) ValidateProbeClient {
	return ValidateProbeClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
