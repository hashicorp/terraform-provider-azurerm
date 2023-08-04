// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/claims"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/graph"
)

type ResourceManagerAccount struct {
	Environment environments.Environment

	ClientId       string
	ObjectId       string
	SubscriptionId string
	TenantId       string

	AuthenticatedAsAServicePrincipal bool
	SkipResourceProviderRegistration bool

	// TODO: delete these when no longer needed by older clients
	AzureEnvironment azure.Environment
}

func NewResourceManagerAccount(ctx context.Context, config auth.Credentials, subscriptionId string, skipResourceProviderRegistration bool, azureEnvironment azure.Environment) (*ResourceManagerAccount, error) {
	authorizer, err := auth.NewAuthorizerFromCredentials(ctx, config, config.Environment.MicrosoftGraph)
	if err != nil {
		return nil, fmt.Errorf("unable to build authorizer for Microsoft Graph API: %+v", err)
	}

	// Acquire an access token so we can inspect the claims
	token, err := authorizer.Token(ctx, &http.Request{})
	if err != nil {
		return nil, fmt.Errorf("could not acquire access token to parse claims: %+v", err)
	}

	claims, err := claims.ParseClaims(token)
	if err != nil {
		return nil, fmt.Errorf("parsing claims from access token: %+v", err)
	}

	authenticatedAsServicePrincipal := true
	if strings.Contains(strings.ToLower(claims.Scopes), "openid") {
		authenticatedAsServicePrincipal = false
	}

	clientId := claims.AppId
	if clientId == "" {
		log.Printf("[DEBUG] Using user-supplied ClientID because the `appid` claim was missing from the access token")
		clientId = config.ClientID
	}

	objectId := claims.ObjectId
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

	tenantId := claims.TenantId
	if tenantId == "" {
		log.Printf("[DEBUG] Using user-supplied TenantID because the `tid` claim was missing from the access token")
		tenantId = config.TenantID
	}

	// Finally, defer to Azure CLI to obtain tenant ID, subscription ID and client ID when not specified and missing from claims
	realAuthorizer := authorizer
	if cache, ok := authorizer.(*auth.CachedAuthorizer); ok {
		realAuthorizer = cache.Source
	}
	if cli, ok := realAuthorizer.(*auth.AzureCliAuthorizer); ok {
		// Use the tenant ID from Azure CLI when otherwise unknown
		if tenantId == "" {
			if cli.TenantID == "" {
				return nil, fmt.Errorf("azure-cli could not determine tenant ID to use")
			}
			tenantId = cli.TenantID
			log.Printf("[DEBUG] Using tenant ID from Azure CLI: %q", tenantId)
		}

		// Use the subscription ID from Azure CLI when otherwise unknown
		if subscriptionId == "" {
			if cli.DefaultSubscriptionID == "" {
				return nil, fmt.Errorf("azure-cli could not determine subscription ID to use and no subscription was specified")
			}

			subscriptionId = cli.DefaultSubscriptionID
			log.Printf("[DEBUG] Using default subscription ID from Azure CLI: %q", subscriptionId)
		}

		// Use the Azure CLI client ID
		if id, ok := config.Environment.MicrosoftAzureCli.AppId(); ok {
			clientId = *id
			log.Printf("[DEBUG] Using client ID from Azure CLI: %q", clientId)
		}
	}

	if tenantId == "" {
		return nil, fmt.Errorf("unable to configure ResourceManagerAccount: tenant ID could not be determined and was not specified")
	}
	if subscriptionId == "" {
		return nil, fmt.Errorf("unable to configure ResourceManagerAccount: subscription ID could not be determined and was not specified")
	}

	account := ResourceManagerAccount{
		Environment: config.Environment,

		ClientId:       clientId,
		ObjectId:       objectId,
		SubscriptionId: subscriptionId,
		TenantId:       tenantId,

		AuthenticatedAsAServicePrincipal: authenticatedAsServicePrincipal,
		SkipResourceProviderRegistration: skipResourceProviderRegistration,

		// TODO: delete these when no longer needed by older clients
		AzureEnvironment: azureEnvironment,
	}

	return &account, nil
}
