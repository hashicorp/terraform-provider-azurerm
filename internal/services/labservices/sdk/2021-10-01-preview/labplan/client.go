package labplan

import "github.com/Azure/go-autorest/autorest"

type LabPlanClient struct {
	Client  autorest.Client
	baseUri string
}

func NewLabPlanClientWithBaseURI(endpoint string) LabPlanClient {
	return LabPlanClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
