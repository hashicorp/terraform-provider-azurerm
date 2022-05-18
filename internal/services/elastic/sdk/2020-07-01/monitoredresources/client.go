package monitoredresources

import "github.com/Azure/go-autorest/autorest"

type MonitoredResourcesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewMonitoredResourcesClientWithBaseURI(endpoint string) MonitoredResourcesClient {
	return MonitoredResourcesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
