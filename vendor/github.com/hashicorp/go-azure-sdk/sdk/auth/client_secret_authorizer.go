// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"golang.org/x/oauth2"
)

type ClientSecretAuthorizerOptions struct {
	// Environment is the Azure environment/cloud being targeted
	Environment environments.Environment

	// Api describes the Azure API being used
	Api environments.Api

	// TenantId is the tenant to authenticate against
	TenantId string

	// AuxTenantIds lists additional tenants to authenticate against, currently only
	// used for Resource Manager when auxiliary tenants are needed.
	// e.g. https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/authenticate-multi-tenant
	AuxTenantIds []string

	// ClientId is the client ID used when authenticating
	ClientId string

	// ClientSecret is the client secret used when authenticating
	ClientSecret string
}

// NewClientSecretAuthorizer returns an authorizer which uses client secret authentication.
func NewClientSecretAuthorizer(ctx context.Context, options ClientSecretAuthorizerOptions) (Authorizer, error) {
	scope, err := environments.Scope(options.Api)
	if err != nil {
		return nil, fmt.Errorf("determining scope for %q: %+v", options.Api.Name(), err)
	}

	conf := clientCredentialsConfig{
		Environment:        options.Environment,
		TenantID:           options.TenantId,
		AuxiliaryTenantIDs: options.AuxTenantIds,
		ClientID:           options.ClientId,
		ClientSecret:       options.ClientSecret,
		Scopes: []string{
			*scope,
		},
	}

	return conf.TokenSource(ctx, clientCredentialsSecretType)
}

var _ Authorizer = &ClientSecretAuthorizer{}

type ClientSecretAuthorizer struct {
	conf *clientCredentialsConfig
}

func (a *ClientSecretAuthorizer) Token(ctx context.Context, _ *http.Request) (*oauth2.Token, error) {
	if a.conf == nil {
		return nil, fmt.Errorf("could not request token: conf is nil")
	}

	v := url.Values{
		"client_id":     {a.conf.ClientID},
		"client_secret": {a.conf.ClientSecret},
		"grant_type":    {"client_credentials"},
		// NOTE: at this time we only support v2 (MSAL) Tokens since v1 (ADAL) is EOL.
		"scope": []string{
			strings.Join(a.conf.Scopes, " "),
		},
	}

	tokenUrl := a.conf.TokenURL
	if tokenUrl == "" {
		if a.conf.Environment.Authorization == nil {
			return nil, fmt.Errorf("no `authorization` configuration was found for this environment")
		}
		tokenUrl = tokenEndpoint(*a.conf.Environment.Authorization, a.conf.TenantID)
	}

	return clientCredentialsToken(ctx, tokenUrl, &v)
}

// AuxiliaryTokens returns additional tokens for auxiliary tenant IDs, for use in multi-tenant scenarios
func (a *ClientSecretAuthorizer) AuxiliaryTokens(ctx context.Context, _ *http.Request) ([]*oauth2.Token, error) {
	if a.conf == nil {
		return nil, fmt.Errorf("could not request token: conf is nil")
	}

	tokens := make([]*oauth2.Token, 0)

	if len(a.conf.AuxiliaryTenantIDs) == 0 {
		return tokens, nil
	}

	for _, tenantId := range a.conf.AuxiliaryTenantIDs {
		v := url.Values{
			"client_id":     {a.conf.ClientID},
			"client_secret": {a.conf.ClientSecret},
			"grant_type":    {"client_credentials"},
			// NOTE: at this time we only support v2 (MSAL) Tokens since v1 (ADAL) is EOL.
			"scope": []string{
				// TODO: given the Request, could we use a dynamic scope?
				strings.Join(a.conf.Scopes, " "),
			},
		}

		tokenUrl := a.conf.TokenURL
		if tokenUrl == "" {
			if a.conf.Environment.Authorization == nil {
				return nil, fmt.Errorf("no `authorization` configuration was found for this environment")
			}
			tokenUrl = tokenEndpoint(*a.conf.Environment.Authorization, tenantId)
		}

		token, err := clientCredentialsToken(ctx, tokenUrl, &v)
		if err != nil {
			return tokens, err
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}
