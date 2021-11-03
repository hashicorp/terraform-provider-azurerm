package webapplicationfirewallmanagedrulesets

import "github.com/Azure/go-autorest/autorest"

type WebApplicationFirewallManagedRuleSetsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewWebApplicationFirewallManagedRuleSetsClientWithBaseURI(endpoint string) WebApplicationFirewallManagedRuleSetsClient {
	return WebApplicationFirewallManagedRuleSetsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
