package privateendpointconnection

import "github.com/Azure/go-autorest/autorest"

type PrivateEndpointConnectionClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPrivateEndpointConnectionClientWithBaseURI(endpoint string) PrivateEndpointConnectionClient {
	return PrivateEndpointConnectionClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
