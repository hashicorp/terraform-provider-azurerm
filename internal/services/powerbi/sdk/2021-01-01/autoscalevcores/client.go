package autoscalevcores

import "github.com/Azure/go-autorest/autorest"

type AutoScaleVCoresClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAutoScaleVCoresClientWithBaseURI(endpoint string) AutoScaleVCoresClient {
	return AutoScaleVCoresClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
