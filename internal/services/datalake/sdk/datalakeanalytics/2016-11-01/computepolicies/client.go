package computepolicies

import "github.com/Azure/go-autorest/autorest"

type ComputePoliciesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewComputePoliciesClientWithBaseURI(endpoint string) ComputePoliciesClient {
	return ComputePoliciesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
