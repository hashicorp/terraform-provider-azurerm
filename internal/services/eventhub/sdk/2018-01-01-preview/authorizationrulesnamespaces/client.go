package authorizationrulesnamespaces

import "github.com/Azure/go-autorest/autorest"

type AuthorizationRulesNamespacesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAuthorizationRulesNamespacesClientWithBaseURI(endpoint string) AuthorizationRulesNamespacesClient {
	return AuthorizationRulesNamespacesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
