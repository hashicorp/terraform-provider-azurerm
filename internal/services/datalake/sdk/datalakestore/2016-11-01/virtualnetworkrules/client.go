package virtualnetworkrules

import "github.com/Azure/go-autorest/autorest"

type VirtualNetworkRulesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewVirtualNetworkRulesClientWithBaseURI(endpoint string) VirtualNetworkRulesClient {
	return VirtualNetworkRulesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
