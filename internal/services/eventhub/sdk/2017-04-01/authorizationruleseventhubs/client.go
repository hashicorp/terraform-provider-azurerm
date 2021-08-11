package authorizationruleseventhubs

import "github.com/Azure/go-autorest/autorest"

type AuthorizationRulesEventHubsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAuthorizationRulesEventHubsClientWithBaseURI(endpoint string) AuthorizationRulesEventHubsClient {
	return AuthorizationRulesEventHubsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
