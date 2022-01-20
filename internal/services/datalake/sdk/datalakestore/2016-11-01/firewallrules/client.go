package firewallrules

import "github.com/Azure/go-autorest/autorest"

type FirewallRulesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewFirewallRulesClientWithBaseURI(endpoint string) FirewallRulesClient {
	return FirewallRulesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
