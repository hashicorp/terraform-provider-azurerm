package dimensions

import "github.com/Azure/go-autorest/autorest"

type DimensionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDimensionsClientWithBaseURI(endpoint string) DimensionsClient {
	return DimensionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
