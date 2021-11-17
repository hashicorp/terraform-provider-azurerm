package alerts

import "github.com/Azure/go-autorest/autorest"

type AlertsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAlertsClientWithBaseURI(endpoint string) AlertsClient {
	return AlertsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
