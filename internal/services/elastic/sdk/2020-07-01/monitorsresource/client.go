package monitorsresource

import "github.com/Azure/go-autorest/autorest"

type MonitorsResourceClient struct {
	Client  autorest.Client
	baseUri string
}

func NewMonitorsResourceClientWithBaseURI(endpoint string) MonitorsResourceClient {
	return MonitorsResourceClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
