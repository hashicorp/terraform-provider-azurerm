// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type OIDCAuthorizerOptions struct {
	// Environment is the Azure environment/cloud being targeted
	Environment environments.Environment

	// Api describes the Azure API being used
	Api environments.Api

	// TenantId is the tenant to authenticate against
	TenantId string

	// AuxiliaryTenantIds lists additional tenants to authenticate against, currently only
	// used for Resource Manager when auxiliary tenants are needed.
	// e.g. https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/authenticate-multi-tenant
	AuxiliaryTenantIds []string

	// ClientId is the client ID used when authenticating
	ClientId string

	// FederatedAssertion is the client assertion dispensed by the OIDC provider used to verify identity during authentication
	FederatedAssertion string
}

// NewOIDCAuthorizer returns an authorizer which uses OIDC authentication (federated client credentials)
func NewOIDCAuthorizer(ctx context.Context, options OIDCAuthorizerOptions) (Authorizer, error) {
	scope, err := environments.Scope(options.Api)
	if err != nil {
		return nil, fmt.Errorf("determining scope for %q: %+v", options.Api.Name(), err)
	}

	conf := clientCredentialsConfig{
		Environment:        options.Environment,
		TenantID:           options.TenantId,
		AuxiliaryTenantIDs: options.AuxiliaryTenantIds,
		ClientID:           options.ClientId,
		FederatedAssertion: options.FederatedAssertion,
		Scopes: []string{
			*scope,
		},
	}

	return conf.TokenSource(ctx, clientCredentialsAssertionType)
}
