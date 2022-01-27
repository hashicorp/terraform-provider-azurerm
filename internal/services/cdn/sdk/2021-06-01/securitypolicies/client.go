package securitypolicies

import "github.com/Azure/go-autorest/autorest"

type SecurityPoliciesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSecurityPoliciesClientWithBaseURI(endpoint string) SecurityPoliciesClient {
	return SecurityPoliciesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
