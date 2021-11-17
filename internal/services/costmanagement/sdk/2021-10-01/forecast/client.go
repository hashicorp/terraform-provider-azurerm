package forecast

import "github.com/Azure/go-autorest/autorest"

type ForecastClient struct {
	Client  autorest.Client
	baseUri string
}

func NewForecastClientWithBaseURI(endpoint string) ForecastClient {
	return ForecastClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
