package vmingestiondetails

import "github.com/Azure/go-autorest/autorest"

type VMIngestionDetailsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewVMIngestionDetailsClientWithBaseURI(endpoint string) VMIngestionDetailsClient {
	return VMIngestionDetailsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
