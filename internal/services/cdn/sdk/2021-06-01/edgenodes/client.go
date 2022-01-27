package edgenodes

import "github.com/Azure/go-autorest/autorest"

type EdgenodesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewEdgenodesClientWithBaseURI(endpoint string) EdgenodesClient {
	return EdgenodesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
