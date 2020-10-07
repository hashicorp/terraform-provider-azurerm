package clients

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
)

type ResourceManagerAccount struct {
	AuthenticatedAsAServicePrincipal bool
	ClientId                         string
	Environment                      azure.Environment
	ObjectId                         string
	SubscriptionId                   string
	TenantId                         string
}

func NewResourceManagerAccount(ctx context.Context, config authentication.Config, env azure.Environment) (*ResourceManagerAccount, error) {
	objectId := ""

	// TODO remove this when we confirm that MSI no longer returns nil with getAuthenticatedObjectID
	if getAuthenticatedObjectID := config.GetAuthenticatedObjectID; getAuthenticatedObjectID != nil {
		v, err := getAuthenticatedObjectID(ctx)
		if err != nil {
			return nil, fmt.Errorf("Error getting authenticated object ID: %v", err)
		}
		objectId = v
	}

	account := ResourceManagerAccount{
		AuthenticatedAsAServicePrincipal: config.AuthenticatedAsAServicePrincipal,
		ClientId:                         config.ClientID,
		Environment:                      env,
		ObjectId:                         objectId,
		TenantId:                         config.TenantID,
		SubscriptionId:                   config.SubscriptionID,
	}
	return &account, nil
}
