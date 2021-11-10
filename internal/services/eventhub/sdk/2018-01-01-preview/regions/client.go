package regions

import "github.com/Azure/go-autorest/autorest"

type RegionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRegionsClientWithBaseURI(endpoint string) RegionsClient {
	return RegionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
