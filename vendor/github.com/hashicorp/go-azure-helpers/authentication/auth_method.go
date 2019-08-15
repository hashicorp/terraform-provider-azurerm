package authentication

import (
	"github.com/Azure/go-autorest/autorest"
)

type authMethod interface {
	build(b Builder) (authMethod, error)

	isApplicable(b Builder) bool

	getAuthorizationToken(sender autorest.Sender, oauthConfig *MultiOAuth, endpoint string) (autorest.Authorizer, error)

	name() string

	populateConfig(c *Config) error

	validate() error
}
