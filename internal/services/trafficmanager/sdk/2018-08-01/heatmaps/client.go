package heatmaps

import "github.com/Azure/go-autorest/autorest"

type HeatMapsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewHeatMapsClientWithBaseURI(endpoint string) HeatMapsClient {
	return HeatMapsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
