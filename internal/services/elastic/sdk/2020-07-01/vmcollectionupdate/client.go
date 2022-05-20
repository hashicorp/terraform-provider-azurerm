package vmcollectionupdate

import "github.com/Azure/go-autorest/autorest"

type VMCollectionUpdateClient struct {
	Client  autorest.Client
	baseUri string
}

func NewVMCollectionUpdateClientWithBaseURI(endpoint string) VMCollectionUpdateClient {
	return VMCollectionUpdateClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
