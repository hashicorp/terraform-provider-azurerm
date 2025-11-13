// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/claims"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/graph"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders"
)

type ResourceManagerAccount struct {
	Environment environments.Environment

	ClientId       string
	ObjectId       string
	SubscriptionId string
	TenantId       string

	AuthenticatedAsAServicePrincipal bool
	RegisteredResourceProviders      resourceproviders.ResourceProviders
}

func NewResourceManagerAccount(ctx context.Context, config auth.Credentials, subscriptionId string, registeredResourceProviders resourceproviders.ResourceProviders) (*ResourceManagerAccount, error) {
	authorizer, err := auth.NewAuthorizerFromCredentials(ctx, config, config.Environment.MicrosoftGraph)
	if err != nil {
		return nil, fmt.Errorf("unable to build authorizer for Microsoft Graph API: %+v", err)
	}

	// Acquire an access token so we can inspect the claims
	token, err := authorizer.Token(ctx, &http.Request{})
	if err != nil {
		return nil, fmt.Errorf("could not acquire access token to parse claims: %+v", err)
	}

	tokenClaims, err := claims.ParseClaims(token)
	if err != nil {
		return nil, fmt.Errorf("parsing claims from access token: %+v", err)
	}

	authenticatedAsServicePrincipal := !strings.Contains(strings.ToLower(tokenClaims.Scopes), "openid")

	clientId := tokenClaims.AppId
	if clientId == "" {
		log.Printf("[DEBUG] Using user-supplied ClientID because the `appid` claim was missing from the access token")
		clientId = config.ClientID
	}

	objectId := tokenClaims.ObjectId
	if objectId == "" {
		if authenticatedAsServicePrincipal {
			log.Printf("[DEBUG] Querying Microsoft Graph to discover authenticated service principal object ID because the `oid` claim was missing from the access token")
			id, err := graph.ServicePrincipalObjectID(ctx, authorizer, config.Environment, config.ClientID)
			if err != nil {
				return nil, fmt.Errorf("attempting to discover object ID for authenticated service principal with client ID %q: %+v", config.ClientID, err)
			}

			objectId = *id
		} else {
			log.Printf("[DEBUG] Querying Microsoft Graph to discover authenticated user principal object ID because the `oid` claim was missing from the access token")
			id, err := graph.UserPrincipalObjectID(ctx, authorizer, config.Environment)
			if err != nil {
				return nil, fmt.Errorf("attempting to discover object ID for authenticated user principal: %+v", err)
			}

			objectId = *id
		}
	}

	tenantId := tokenClaims.TenantId
	if tenantId == "" {
		log.Printf("[DEBUG] Using user-supplied TenantID because the `tid` claim was missing from the access token")
		tenantId = config.TenantID
	}

	// Finally, defer to Azure CLI to obtain tenant ID and client ID when not specified and missing from claims
	realAuthorizer := authorizer
	if cache, ok := authorizer.(*auth.CachedAuthorizer); ok {
		realAuthorizer = cache.Source
	}
	if cli, ok := realAuthorizer.(*auth.AzureCliAuthorizer); ok {
		// Use the tenant ID from Azure CLI when otherwise unknown
		if tenantId == "" {
			if cli.TenantID == "" {
				return nil, errors.New("azure-cli could not determine tenant ID to use")
			}
			tenantId = cli.TenantID
			log.Printf("[DEBUG] Using tenant ID from Azure CLI: %q", tenantId)
		}

		// Use the Azure CLI client ID
		if id, ok := config.Environment.MicrosoftAzureCli.AppId(); ok {
			clientId = *id
			log.Printf("[DEBUG] Using client ID from Azure CLI: %q", clientId)
		}
	}

	// We'll permit the provider to proceed with an unknown client ID since it only affects a small number of use cases when authenticating as a user
	if tenantId == "" {
		return nil, errors.New("unable to configure ResourceManagerAccount: tenant ID could not be determined and was not specified")
	}
	if subscriptionId == "" {
		return nil, errors.New("unable to configure ResourceManagerAccount: subscription ID could not be determined and was not specified NOOOOOOOOOB!!!!")
	}

	account := ResourceManagerAccount{
		Environment: config.Environment,

		ClientId:       clientId,
		ObjectId:       objectId,
		SubscriptionId: subscriptionId,
		TenantId:       tenantId,

		AuthenticatedAsAServicePrincipal: authenticatedAsServicePrincipal,
		RegisteredResourceProviders:      registeredResourceProviders,
	}

	return &account, nil
}
