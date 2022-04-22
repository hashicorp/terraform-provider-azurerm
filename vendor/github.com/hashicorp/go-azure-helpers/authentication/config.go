package authentication

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	authWrapper "github.com/manicminer/hamilton-autorest/auth"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
)

// Config is the configuration structure used to instantiate a
// new Azure management client.
type Config struct {
	ClientID           string
	SubscriptionID     string
	TenantID           string
	AuxiliaryTenantIDs []string
	Environment        string
	MetadataHost       string

	GetAuthenticatedObjectID         func(context.Context) (*string, error)
	AuthenticatedAsAServicePrincipal bool

	// A Custom Resource Manager Endpoint
	// at this time this should only be applicable for Azure Stack.
	CustomResourceManagerEndpoint string

	// Beta opt-in for Microsoft Graph
	UseMicrosoftGraph bool

	authMethod authMethod
}

type OAuthConfig struct {
	OAuth            *adal.OAuthConfig
	MultiTenantOauth *adal.MultiTenantOAuthConfig
}

// GetAuthorizationToken returns an authorization token for the authentication method defined in the Config
func (c Config) GetOAuthConfig(activeDirectoryEndpoint string) (*adal.OAuthConfig, error) {
	log.Printf("Getting OAuth config for endpoint %s with  tenant %s", activeDirectoryEndpoint, c.TenantID)

	// fix for ADFS environments, if the login endpoint ends in `/adfs` it's an adfs environment
	// the login endpoint ends up residing in `ActiveDirectoryEndpoint`
	oAuthTenant := c.TenantID
	if strings.HasSuffix(strings.ToLower(activeDirectoryEndpoint), "/adfs") {
		log.Printf("[DEBUG] ADFS environment detected - overriding Tenant ID to `adfs`!")
		oAuthTenant = "adfs"
	}

	oauth, err := adal.NewOAuthConfig(activeDirectoryEndpoint, oAuthTenant)
	if err != nil {
		return nil, err
	}

	// OAuthConfigForTenant returns a pointer, which can be nil.
	if oauth == nil {
		return nil, fmt.Errorf("unable to configure OAuthConfig for tenant %s", c.TenantID)
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
		return nil, fmt.Errorf("unable to configure OAuthConfig for tenant %s (auxiliary tenants %v)", c.TenantID, c.AuxiliaryTenantIDs)
	}

	return &oauth, nil
}

// BuildOAuthConfig builds the authorization configuration for the specified Active Directory Endpoint
func (c Config) BuildOAuthConfig(activeDirectoryEndpoint string) (*OAuthConfig, error) {
	multiAuth := OAuthConfig{}
	var err error

	multiAuth.OAuth, err = c.GetOAuthConfig(activeDirectoryEndpoint)
	if err != nil {
		return nil, err
	}

	if len(c.AuxiliaryTenantIDs) > 0 {
		multiAuth.MultiTenantOauth, err = c.GetMultiTenantOAuthConfig(activeDirectoryEndpoint)
		if err != nil {
			return nil, err
		}
	}

	return &multiAuth, nil
}

// ADALBearerAuthorizerCallback returns a BearerAuthorizer valid only for the Primary Tenant
// this signs a request using the AccessToken returned from the primary Resource Manager authorizer
func (c Config) ADALBearerAuthorizerCallback(ctx context.Context, sender autorest.Sender, oauthConfig *OAuthConfig) *autorest.BearerAuthorizerCallback {
	return autorest.NewBearerAuthorizerCallback(sender, func(tenantID, resource string) (*autorest.BearerAuthorizer, error) {
		// a BearerAuthorizer is only valid for the primary tenant
		newAuthConfig := &OAuthConfig{
			OAuth: oauthConfig.OAuth,
		}

		storageSpt, err := c.GetADALToken(ctx, sender, newAuthConfig, resource)
		if err != nil {
			return nil, err
		}

		cast, ok := storageSpt.(*autorest.BearerAuthorizer)
		if !ok {
			return nil, fmt.Errorf("converting %+v to a BearerAuthorizer", storageSpt)
		}

		return cast, nil
	})
}

// MSALBearerAuthorizerCallback returns a BearerAuthorizer valid only for the Primary Tenant
// this signs a request using the AccessToken returned from the primary Resource Manager authorizer
func (c Config) MSALBearerAuthorizerCallback(ctx context.Context, api environments.Api, sender autorest.Sender, oauthConfig *OAuthConfig, endpoint string) *autorest.BearerAuthorizerCallback {
	authorizer, err := c.GetMSALToken(ctx, api, sender, oauthConfig, endpoint)
	if err != nil {
		return autorest.NewBearerAuthorizerCallback(nil, func(_, _ string) (*autorest.BearerAuthorizer, error) {
			return nil, fmt.Errorf("failed to acquire MSAL token for %s", api.Endpoint)
		})
	}

	// For compatibility with Azure CLI which still uses ADAL, first check if we got an autorest.BearerAuthorizer
	if cast, ok := authorizer.(*autorest.BearerAuthorizer); ok {
		return autorest.NewBearerAuthorizerCallback(sender, func(_, _ string) (*autorest.BearerAuthorizer, error) {
			return cast, nil
		})
	}

	cast, ok := authorizer.(auth.Authorizer)
	if !ok {
		return autorest.NewBearerAuthorizerCallback(nil, func(_, _ string) (*autorest.BearerAuthorizer, error) {
			return nil, fmt.Errorf("authorizer was not an auth.Authorizer for %s", api.Endpoint)
		})
	}

	return (&authWrapper.Authorizer{Authorizer: cast}).BearerAuthorizerCallback()
}

// GetADALToken returns an autorest.Authorizer using an ADAL token via the authentication method defined in the Config
func (c Config) GetADALToken(ctx context.Context, sender autorest.Sender, oauth *OAuthConfig, endpoint string) (autorest.Authorizer, error) {
	return c.authMethod.getADALToken(ctx, sender, oauth, endpoint)
}

// GetMSALToken returns an autorest.Authorizer using an MSAL token via the authentication method defined in the Config
func (c Config) GetMSALToken(ctx context.Context, api environments.Api, sender autorest.Sender, oauth *OAuthConfig, endpoint string) (autorest.Authorizer, error) {
	return c.authMethod.getMSALToken(ctx, api, sender, oauth, endpoint)
}
