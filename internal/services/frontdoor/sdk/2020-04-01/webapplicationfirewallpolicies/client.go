package webapplicationfirewallpolicies

import "github.com/Azure/go-autorest/autorest"

type WebApplicationFirewallPoliciesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewWebApplicationFirewallPoliciesClientWithBaseURI(endpoint string) WebApplicationFirewallPoliciesClient {
	return WebApplicationFirewallPoliciesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
