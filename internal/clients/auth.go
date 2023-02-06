package clients

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/claims"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
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
	AzureEnvironment        azure.Environment
	ResourceManagerEndpoint string
}

func NewResourceManagerAccount(ctx context.Context, authorizer auth.Authorizer, config auth.Credentials, subscriptionId string, skipResourceProviderRegistration bool, azureEnvironment azure.Environment, resourceManagerEndpoint string) (*ResourceManagerAccount, error) {
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
	if strings.Contains(strings.ToLower(claims.Scopes), "user_impersonation") {
		authenticatedAsServicePrincipal = false
	}

	account := ResourceManagerAccount{
		Environment: config.Environment,

		ClientId:       claims.AppId,
		ObjectId:       claims.ObjectId,
		SubscriptionId: subscriptionId,
		TenantId:       claims.TenantId,

		AuthenticatedAsAServicePrincipal: authenticatedAsServicePrincipal,
		SkipResourceProviderRegistration: skipResourceProviderRegistration,

		// TODO: delete these when no longer needed by older clients
		AzureEnvironment:        azureEnvironment,
		ResourceManagerEndpoint: resourceManagerEndpoint,
	}

	return &account, nil
}
