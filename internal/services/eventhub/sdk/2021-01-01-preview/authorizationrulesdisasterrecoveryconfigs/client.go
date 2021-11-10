package authorizationrulesdisasterrecoveryconfigs

import "github.com/Azure/go-autorest/autorest"

type AuthorizationRulesDisasterRecoveryConfigsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewAuthorizationRulesDisasterRecoveryConfigsClientWithBaseURI(endpoint string) AuthorizationRulesDisasterRecoveryConfigsClient {
	return AuthorizationRulesDisasterRecoveryConfigsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
