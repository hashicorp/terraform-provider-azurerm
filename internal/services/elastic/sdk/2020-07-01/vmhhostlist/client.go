package vmhhostlist

import "github.com/Azure/go-autorest/autorest"

type VMHHostListClient struct {
	Client  autorest.Client
	baseUri string
}

func NewVMHHostListClientWithBaseURI(endpoint string) VMHHostListClient {
	return VMHHostListClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
