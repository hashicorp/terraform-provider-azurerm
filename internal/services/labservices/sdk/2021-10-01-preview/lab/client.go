package lab

import "github.com/Azure/go-autorest/autorest"

type LabClient struct {
	Client  autorest.Client
	baseUri string
}

func NewLabClientWithBaseURI(endpoint string) LabClient {
	return LabClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
