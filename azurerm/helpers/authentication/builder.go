package authentication

import (
	"fmt"
	"log"
)

// Builder supports all of the possible Authentication values and feature toggles
// required to build a working Config for Authentication purposes.
type Builder struct {
	// Core
	ClientID       string
	SubscriptionID string
	TenantID       string
	Environment    string

	// Azure CLI Parsing / CloudShell Auth
	SupportsAzureCliCloudShellParsing bool

	// Managed Service Identity Auth
	SupportsManagedServiceIdentity bool
	MsiEndpoint                    string

	// Service Principal (Client Cert) Auth
	SupportsClientCertAuth bool
	ClientCertPath         string
	ClientCertPassword     string

	// Service Principal (Client Secret) Auth
	SupportsClientSecretAuth bool
	ClientSecret             string
}

// Build takes the configuration from the Builder and builds up a validated Config
// for authenticating with Azure
func (b Builder) Build() (*Config, error) {
	config := Config{
		ClientID:       b.ClientID,
		SubscriptionID: b.SubscriptionID,
		TenantID:       b.TenantID,
		Environment:    b.Environment,
	}

	if b.SupportsClientCertAuth && b.ClientCertPath != "" {
		log.Printf("[DEBUG] Using Service Principal / Client Certificate for Authentication")
		config.AuthenticatedAsAServicePrincipal = true

		method, err := newServicePrincipalClientCertificateAuth(b)
		if err != nil {
			return nil, err
		}

		config.authMethod = method
		return config.validate()
	}

	if b.SupportsClientSecretAuth && b.ClientSecret != "" {
		log.Printf("[DEBUG] Using Service Principal / Client Secret for Authentication")
		config.AuthenticatedAsAServicePrincipal = true

		method, err := newServicePrincipalClientSecretAuth(b)
		if err != nil {
			return nil, err
		}

		config.authMethod = method
		return config.validate()
	}

	if b.SupportsManagedServiceIdentity {
		log.Printf("[DEBUG] Using Managed Service Identity for Authentication")
		method, err := newManagedServiceIdentityAuth(b)
		if err != nil {
			return nil, err
		}
		config.authMethod = method
		return config.validate()
	}

	if b.SupportsAzureCliCloudShellParsing {
		log.Printf("[DEBUG] Parsing credentials from the Azure CLI for Authentication")

		method, err := newAzureCliParsingAuth(b)
		if err != nil {
			return nil, err
		}

		// as credentials are parsed from the Azure CLI's Profile we actually need to
		// obtain the ClientId, Environment, Subscription ID & TenantID here
		if cliAuth, ok := method.(azureCliParsingAuth); ok && method != nil {
			config.ClientID = cliAuth.profile.clientId
			config.Environment = cliAuth.profile.environment
			config.SubscriptionID = cliAuth.profile.subscriptionId
			config.TenantID = cliAuth.profile.tenantId
		}

		config.authMethod = method
		return config.validate()
	}

	return nil, fmt.Errorf("No supported authentication methods were found!")
}
