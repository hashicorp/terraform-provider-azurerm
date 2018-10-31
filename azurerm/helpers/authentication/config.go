package authentication

import (
	"github.com/Azure/go-autorest/autorest/adal"
)

// Config is the configuration structure used to instantiate a
// new Azure management client.
type Config struct {
	// Core
	ClientID       string
	SubscriptionID string
	TenantID       string
	Environment    string

	// temporarily public feature flags
	AuthenticatedAsAServicePrincipal bool

	// to be removed
	accessToken              *adal.Token
	usingCloudShell          bool
	authMethod               authMethod
	SkipProviderRegistration bool
}

func (c Config) validate() (*Config, error) {
	err := c.authMethod.validate()
	if err != nil {
		return nil, err
	}

	return &c, nil
}
