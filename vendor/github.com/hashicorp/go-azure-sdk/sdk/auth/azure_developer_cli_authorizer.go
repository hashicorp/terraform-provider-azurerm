// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/go-azure-sdk/sdk/internal/azuredevelopercli"
	"golang.org/x/oauth2"
)

type AzureDeveloperCliAuthorizerOptions struct {
	// TenantId is the tenant to authenticate against
	TenantId string

	// AuxTenantIds lists additional tenants to authenticate against, currently only
	// used for Resource Manager when auxiliary tenants are needed.
	// e.g. https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/authenticate-multi-tenant
	AuxTenantIds []string

	// Api describes the Azure API being used
	Api environments.Api
}

// NewAzureDeveloperCliAuthorizer returns an Authorizer which authenticates using the Azure Developer CLI.
func NewAzureDeveloperCliAuthorizer(ctx context.Context, options AzureDeveloperCliAuthorizerOptions) (Authorizer, error) {
	conf, err := newAzureDeveloperCliConfig(options.Api, options.TenantId, options.AuxTenantIds)
	if err != nil {
		return nil, err
	}
	return conf.TokenSource(ctx)
}

var _ Authorizer = &AzureDeveloperCliAuthorizer{}

// AzureDeveloperCliAuthorizer is an Authorizer which supports the Azure Developer CLI (azd).
type AzureDeveloperCliAuthorizer struct {
	// TenantID is the specified tenant ID, or the auto-detected tenant ID if none was specified
	TenantID string

	conf *azureDeveloperCliConfig
}

// Token returns an access token using the Azure Developer CLI as an authentication mechanism.
func (a *AzureDeveloperCliAuthorizer) Token(_ context.Context, _ *http.Request) (*oauth2.Token, error) {
	if a.conf == nil {
		return nil, fmt.Errorf("could not request token: conf is nil")
	}

	azArgs := []string{"auth", "token"}

	scope, err := environments.Scope(a.conf.Api)
	if err != nil {
		return nil, fmt.Errorf("determining scope for %q: %+v", a.conf.Api.Name(), err)
	}
	azArgs = append(azArgs, "--scope", *scope)

	// Try to detect if we're running in Cloud Shell
	if cloudShell := os.Getenv("AZUREPS_HOST_ENVIRONMENT"); !strings.HasPrefix(cloudShell, "cloud-shell/") {
		// Seemingly not, so we'll append the tenant ID to the az args
		azArgs = append(azArgs, "--tenant-id", a.conf.TenantID)
	}

	var token azureDeveloperCliToken
	if err := azuredevelopercli.JSONUnmarshalAzdCmd(&token, azArgs...); err != nil {
		return nil, err
	}

	var expiry time.Time
	if token.ExpiresOn != "" {
		if expiry, err = time.ParseInLocation("2006-01-02T15:04:05Z", token.ExpiresOn, time.Local); err != nil {
			return nil, fmt.Errorf("internal-error: parsing expiresOn value for azd-cli auth token")
		}
	}

	return &oauth2.Token{
		AccessToken: token.AccessToken,
		Expiry:      expiry,
		TokenType:   "Bearer",
	}, nil
}

// AuxiliaryTokens returns additional tokens for auxiliary tenant IDs, for use in multi-tenant scenarios
func (a *AzureDeveloperCliAuthorizer) AuxiliaryTokens(_ context.Context, _ *http.Request) ([]*oauth2.Token, error) {
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

	azArgs := []string{"auth", "token"}

	scope, err := environments.Scope(a.conf.Api)
	if err != nil {
		return nil, fmt.Errorf("determining scope for %q: %+v", a.conf.Api.Name(), err)
	}
	azArgs = append(azArgs, "--scope", *scope)

	tokens := make([]*oauth2.Token, 0)
	for _, tenantId := range a.conf.AuxiliaryTenantIDs {
		argsWithTenant := append(azArgs, "--tenant-id", tenantId)

		var token azureDeveloperCliToken
		if err := azuredevelopercli.JSONUnmarshalAzdCmd(&token, argsWithTenant...); err != nil {
			return nil, err
		}

		tokens = append(tokens, &oauth2.Token{
			AccessToken: token.AccessToken,
			TokenType:   "Bearer",
		})
	}

	return tokens, nil
}

// azureDeveloperCliConfig configures an AzureDeveloperCliAuthorizer.
type azureDeveloperCliConfig struct {
	Api environments.Api

	// TenantID is the required tenant ID for the primary token
	TenantID string

	// AuxiliaryTenantIDs is an optional list of tenant IDs for which to obtain additional tokens
	AuxiliaryTenantIDs []string
}

// newAzureDeveloperCliConfig returns a new azureDeveloperCliConfig.
func newAzureDeveloperCliConfig(api environments.Api, tenantId string, auxiliaryTenantIds []string) (*azureDeveloperCliConfig, error) {
	return &azureDeveloperCliConfig{
		Api:                api,
		TenantID:           tenantId,
		AuxiliaryTenantIDs: auxiliaryTenantIds,
	}, nil
}

// TokenSource provides a source for obtaining access tokens using AzureDeveloperCliAuthorizer.
func (c *azureDeveloperCliConfig) TokenSource(ctx context.Context) (Authorizer, error) {
	// Cache access tokens internally to avoid unnecessary `azd` invocations
	return NewCachedAuthorizer(&AzureDeveloperCliAuthorizer{
		TenantID: c.TenantID,
		conf:     c,
	})
}

type azureDeveloperCliToken struct {
	AccessToken string `json:"token"`
	ExpiresOn   string `json:"expiresOn"`
}
