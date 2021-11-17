package usagedetails

import "github.com/Azure/go-autorest/autorest"

type UsageDetailsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewUsageDetailsClientWithBaseURI(endpoint string) UsageDetailsClient {
	return UsageDetailsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
