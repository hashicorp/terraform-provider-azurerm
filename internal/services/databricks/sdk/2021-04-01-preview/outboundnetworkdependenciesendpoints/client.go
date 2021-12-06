package outboundnetworkdependenciesendpoints

import "github.com/Azure/go-autorest/autorest"

type OutboundNetworkDependenciesEndpointsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewOutboundNetworkDependenciesEndpointsClientWithBaseURI(endpoint string) OutboundNetworkDependenciesEndpointsClient {
	return OutboundNetworkDependenciesEndpointsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
