// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func GetSendKeyValue(ctx context.Context, metadata sdk.ResourceMetaData, id webapps.RelayId, sendKeyName string) (*string, error) {
	hybridConnectionsClient := metadata.Client.Relay.HybridConnectionsClient
	connectionId := hybridconnections.NewHybridConnectionAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.HybridConnectionNamespaceName, id.RelayName, sendKeyName)
	keys, err := hybridConnectionsClient.ListKeys(ctx, connectionId)
	if err != nil && !response.WasNotFound(keys.HttpResponse) {
		return nil, fmt.Errorf("listing Send Keys for name %s for %s in %s: %+v", sendKeyName, connectionId, id, err)
	}
	if keys.Model != nil && keys.Model.PrimaryKey != nil {
		return keys.Model.PrimaryKey, nil
	}

	relayNamespaceClient := metadata.Client.Relay.NamespacesClient
	relayConnectionId := namespaces.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.HybridConnectionNamespaceName, sendKeyName)
	relayKeys, err := relayNamespaceClient.ListKeys(ctx, relayConnectionId)
	if err != nil {
		return nil, fmt.Errorf("listing Send Keys for name %s for %s in %s: %+v", sendKeyName, relayConnectionId, id, err)
	}
	if relayKeys.Model == nil || relayKeys.Model.PrimaryKey == nil {
		return nil, fmt.Errorf("reading Send Key Value for %s in %s", relayConnectionId.AuthorizationRuleName, id)
	}
	return relayKeys.Model.PrimaryKey, nil
}
