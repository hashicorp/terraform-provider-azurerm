package connectiongateways

import "github.com/Azure/go-autorest/autorest"

type ConnectionGatewaysClient struct {
	Client  autorest.Client
	baseUri string
}

func NewConnectionGatewaysClientWithBaseURI(endpoint string) ConnectionGatewaysClient {
	return ConnectionGatewaysClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
