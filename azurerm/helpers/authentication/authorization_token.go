package authentication

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
)

// TODO: this can become c.GetAuthorizationToken(oAuthConfig *adal.OAuthConfig, endpoint string)
func GetAuthorizationToken(c *Config, oauthConfig *adal.OAuthConfig, endpoint string) (*autorest.BearerAuthorizer, error) {
	return c.authMethod.getAuthorizationToken(oauthConfig, endpoint)
}
