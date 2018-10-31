package authentication

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

	authMethod authMethod

	// to be removed
	SkipProviderRegistration bool
}

func (c Config) validate() (*Config, error) {
	err := c.authMethod.validate()
	if err != nil {
		return nil, err
	}

	return &c, nil
}
