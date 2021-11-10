package ipfilterrules

import "github.com/Azure/go-autorest/autorest"

type IpFilterRulesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewIpFilterRulesClientWithBaseURI(endpoint string) IpFilterRulesClient {
	return IpFilterRulesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
