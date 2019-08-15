package authentication

import (
	`fmt`
	"log"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
)

// Config is the configuration structure used to instantiate a
// new Azure management client.
type Config struct {
	ClientID                         string
	SubscriptionID                   string
	TenantID                         string
	AuxiliaryTenantIDs               []string
	Environment                      string
	AuthenticatedAsAServicePrincipal bool

	// A Custom Resource Manager Endpoint
	// at this time this should only be applicable for Azure Stack.
	CustomResourceManagerEndpoint string

	authMethod authMethod
}

type MultiOAuth struct {
	OAuth            *adal.OAuthConfig
	MultiTenantOauth *adal.MultiTenantOAuthConfig
}

// GetAuthorizationToken returns an authorization token for the authentication method defined in the Config
func (c Config) GetOAuthConfig(activeDirectoryEndpoint string) (*adal.OAuthConfig, error) {
	log.Printf("Getting OAuth config for endpoint %s with  tenant %s", activeDirectoryEndpoint, c.TenantID)
	oauth, err := adal.NewOAuthConfig(activeDirectoryEndpoint, c.TenantID)
	if err != nil {
		return nil, err
	}

	// OAuthConfigForTenant returns a pointer, which can be nil.
	if oauth == nil {
		return nil, fmt.Errorf("Unable to configure OAuthConfig for tenant %s", c.TenantID)
	}

	return oauth, nil
}

// GetMultiTenantOAuthConfig returns a multi-tenant authorization token for the authentication method defined in the Config
func (c Config) GetMultiTenantOAuthConfig(activeDirectoryEndpoint string) (*adal.MultiTenantOAuthConfig, error) {
	log.Printf("Getting multi OAuth config for endpoint %s with  tenant %s (aux tenants: %v)", activeDirectoryEndpoint, c.TenantID, c.AuxiliaryTenantIDs)
	oauth, err := adal.NewMultiTenantOAuthConfig(activeDirectoryEndpoint, c.TenantID, c.AuxiliaryTenantIDs, adal.OAuthOptions{})
	if err != nil {
		return nil, err
	}

	// OAuthConfigForTenant returns a pointer, which can be nil.
	if oauth == nil {
		return nil, fmt.Errorf("Unable to configure OAuthConfig for tenant %s (auxiliary tenants %v)", c.TenantID, c.AuxiliaryTenantIDs)
	}

	return &oauth, nil
}

func (c Config) GetMultiOAuthConfig(activeDirectoryEndpoint string) (*MultiOAuth, error) {
	if len(c.AuxiliaryTenantIDs) == 0 {
		oauth, err := c.GetOAuthConfig(activeDirectoryEndpoint)
		return &MultiOAuth{OAuth: oauth}, err
	}

	oauth, err := c.GetMultiTenantOAuthConfig(activeDirectoryEndpoint)
	return &MultiOAuth{MultiTenantOauth: oauth}, err
}

// GetAuthorizationToken returns an authorization token for the authentication method defined in the Config
func (c Config) GetAuthorizationToken(sender autorest.Sender, oauth *MultiOAuth, endpoint string) (autorest.Authorizer, error) {
	return c.authMethod.getAuthorizationToken(sender, oauth, endpoint)
}
