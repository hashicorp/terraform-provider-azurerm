// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/go-azure-sdk/sdk/internal/azurecli"
	"golang.org/x/oauth2"
)

type AzureCliAuthorizerOptions struct {
	// Api describes the Azure API being used
	Api environments.Api

	// TenantId is the tenant to authenticate against
	TenantId string

	// AuxTenantIds lists additional tenants to authenticate against, currently only
	// used for Resource Manager when auxiliary tenants are needed.
	// e.g. https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/authenticate-multi-tenant
	AuxTenantIds []string
}

// NewAzureCliAuthorizer returns an Authorizer which authenticates using the Azure CLI.
func NewAzureCliAuthorizer(ctx context.Context, options AzureCliAuthorizerOptions) (Authorizer, error) {
	conf, err := newAzureCliConfig(options.Api, options.TenantId, options.AuxTenantIds)
	if err != nil {
		return nil, err
	}
	return conf.TokenSource(ctx)
}

var _ Authorizer = &AzureCliAuthorizer{}

// AzureCliAuthorizer is an Authorizer which supports the Azure CLI.
type AzureCliAuthorizer struct {
	// TenantID is the specified tenant ID, or the auto-detected tenant ID if none was specified
	TenantID string

	// DefaultSubscriptionID is the default subscription, when detected
	DefaultSubscriptionID string

	conf *azureCliConfig
}

// Token returns an access token using the Azure CLI as an authentication mechanism.
func (a *AzureCliAuthorizer) Token(_ context.Context, _ *http.Request) (*oauth2.Token, error) {
	if a.conf == nil {
		return nil, fmt.Errorf("could not request token: conf is nil")
	}

	azArgs := []string{"account", "get-access-token"}

	// verify that the Azure CLI supports MSAL - ADAL is no longer supported
	err := azurecli.CheckAzVersion(azurecli.MsalVersion, nil)
	if err != nil {
		return nil, fmt.Errorf("checking the version of the Azure CLI: %+v", err)
	}
	scope, err := environments.Scope(a.conf.Api)
	if err != nil {
		return nil, fmt.Errorf("determining scope for %q: %+v", a.conf.Api.Name(), err)
	}
	azArgs = append(azArgs, "--scope", *scope)

	// Try to detect if we're running in Cloud Shell
	if cloudShell := os.Getenv("AZUREPS_HOST_ENVIRONMENT"); !strings.HasPrefix(cloudShell, "cloud-shell/") {
		// Seemingly not, so we'll append the tenant ID to the az args
		azArgs = append(azArgs, "--tenant", a.conf.TenantID)
	}

	var token azureCliToken
	if err := azurecli.JSONUnmarshalAzCmd(&token, azArgs...); err != nil {
		return nil, err
	}

	var expiry time.Time
	if token.ExpiresOn != "" {
		if expiry, err = time.ParseInLocation("2006-01-02 15:04:05.999999", token.ExpiresOn, time.Local); err != nil {
			return nil, fmt.Errorf("internal-error: parsing expiresOn value for az-cli token")
		}
	}

	return &oauth2.Token{
		AccessToken: token.AccessToken,
		Expiry:      expiry,
		TokenType:   token.TokenType,
	}, nil
}

// AuxiliaryTokens returns additional tokens for auxiliary tenant IDs, for use in multi-tenant scenarios
func (a *AzureCliAuthorizer) AuxiliaryTokens(_ context.Context, _ *http.Request) ([]*oauth2.Token, error) {
	if a.conf == nil {
		return nil, fmt.Errorf("could not request token: conf is nil")
	}

	// Return early if no auxiliary tenants are configured
	if len(a.conf.AuxiliaryTenantIDs) == 0 {
		return []*oauth2.Token{}, nil
	}

	// Try to detect if we're running in Cloud Shell
	if cloudShell := os.Getenv("AZUREPS_HOST_ENVIRONMENT"); strings.HasPrefix(cloudShell, "cloud-shell/") {
		return nil, fmt.Errorf("auxiliary tokens not supported in Cloud Shell")
	}

	azArgs := []string{"account", "get-access-token"}

	// verify that the Azure CLI supports MSAL - ADAL is no longer supported
	err := azurecli.CheckAzVersion(AzureCliMsalVersion, nil)
	if err != nil {
		return nil, fmt.Errorf("checking the version of the Azure CLI: %+v", err)
	}
	scope, err := environments.Scope(a.conf.Api)
	if err != nil {
		return nil, fmt.Errorf("determining scope for %q: %+v", a.conf.Api.Name(), err)
	}
	azArgs = append(azArgs, "--scope", *scope)

	tokens := make([]*oauth2.Token, 0)
	for _, tenantId := range a.conf.AuxiliaryTenantIDs {
		argsWithTenant := append(azArgs, "--tenant", tenantId)

		var token azureCliToken
		if err := azurecli.JSONUnmarshalAzCmd(&token, argsWithTenant...); err != nil {
			return nil, err
		}

		tokens = append(tokens, &oauth2.Token{
			AccessToken: token.AccessToken,
			TokenType:   token.TokenType,
		})
	}

	return tokens, nil
}

const (
	AzureCliMinimumVersion   = "2.0.81"
	AzureCliMsalVersion      = "2.30.0"
	AzureCliNextMajorVersion = "3.0.0"
)

// azureCliConfig configures an AzureCliAuthorizer.
type azureCliConfig struct {
	Api environments.Api

	// TenantID is the required tenant ID for the primary token
	TenantID string

	// AuxiliaryTenantIDs is an optional list of tenant IDs for which to obtain additional tokens
	AuxiliaryTenantIDs []string

	// DefaultSubscriptionID is the optional default subscription ID
	DefaultSubscriptionID string
}

// newAzureCliConfig validates the supplied tenant ID and returns a new azureCliConfig.
func newAzureCliConfig(api environments.Api, tenantId string, auxiliaryTenantIds []string) (*azureCliConfig, error) {
	var err error

	// check az-cli version
	nextMajor := azurecli.NextMajorVersion
	if err = azurecli.CheckAzVersion(azurecli.MinimumVersion, &nextMajor); err != nil {
		return nil, err
	}

	// check tenant ID
	tenantId, err = azurecli.CheckTenantID(tenantId)
	if err != nil {
		return nil, err
	}
	if tenantId == "" {
		return nil, errors.New("invalid tenantId or unable to determine tenantId")
	}

	// get the default subscription ID
	subscriptionId, err := azurecli.GetDefaultSubscriptionID()
	if err != nil {
		return nil, err
	}

	return &azureCliConfig{
		Api:                   api,
		TenantID:              tenantId,
		AuxiliaryTenantIDs:    auxiliaryTenantIds,
		DefaultSubscriptionID: subscriptionId,
	}, nil
}

// TokenSource provides a source for obtaining access tokens using AzureCliAuthorizer.
func (c *azureCliConfig) TokenSource(ctx context.Context) (Authorizer, error) {
	// Cache access tokens internally to avoid unnecessary `az` invocations
	return NewCachedAuthorizer(&AzureCliAuthorizer{
		TenantID:              c.TenantID,
		DefaultSubscriptionID: c.DefaultSubscriptionID,
		conf:                  c,
	})
}

type azureCliToken struct {
	AccessToken string `json:"accessToken"`
	ExpiresOn   string `json:"expiresOn"`
	Tenant      string `json:"tenant"`
	TokenType   string `json:"tokenType"`
}
