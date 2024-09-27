// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func GetSendKeyValue(ctx context.Context, metadata sdk.ResourceMetaData, id hybridconnections.HybridConnectionId, sendKeyName string) (*string, error) {
	relayNamespaceClient := metadata.Client.Relay.NamespacesClient
	relayConnectionId := namespaces.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, sendKeyName)
	relayKeys, err := relayNamespaceClient.ListKeys(ctx, relayConnectionId)
	if err != nil && !response.WasNotFound(relayKeys.HttpResponse) {
		return nil, fmt.Errorf("listing Send Keys for name %s for %s in %s: %+v", sendKeyName, relayConnectionId, id, err)
	}
	if relayKeys.Model != nil && relayKeys.Model.PrimaryKey != nil {
		return relayKeys.Model.PrimaryKey, nil
	}

	hybridConnectionsClient := metadata.Client.Relay.HybridConnectionsClient
	connectionId := hybridconnections.NewHybridConnectionAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.HybridConnectionName, sendKeyName)
	keys, err := hybridConnectionsClient.ListKeys(ctx, connectionId)
	if err != nil {
		return nil, fmt.Errorf("listing Send Keys for name %s for %s in %s: %+v", sendKeyName, connectionId, id, err)
	}
	if keys.Model == nil || keys.Model.PrimaryKey == nil {
		return nil, fmt.Errorf("reading Send Key Value for %s in %s", connectionId.AuthorizationRuleName, id)
	}
	return keys.Model.PrimaryKey, nil
}
