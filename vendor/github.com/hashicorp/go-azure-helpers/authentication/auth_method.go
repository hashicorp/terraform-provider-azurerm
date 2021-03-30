package authentication

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/go-autorest/autorest"
)

type authMethodBase interface {

	isApplicable(b Builder) bool

	name() string

	populateConfig(c *Config) error

	validate() error
}

type authMethod interface {
	authMethodBase

	build(b Builder) (authMethod, error)

	getAuthorizationToken(sender autorest.Sender, oauthConfig *OAuthConfig, endpoint string) (autorest.Authorizer, error)
}

type authMethodTrack2 interface {
	authMethodBase

	buildAuthMethodTrack2(b Builder) (authMethodTrack2, error)

	getTokenCredential(endpoint string) (azcore.TokenCredential, error)
}
