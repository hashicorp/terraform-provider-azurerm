package authentication

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
)

type authMethod interface {
	getAuthorizationToken(oauthConfig *adal.OAuthConfig, endpoint string) (*autorest.BearerAuthorizer, error)

	validate() error
}
