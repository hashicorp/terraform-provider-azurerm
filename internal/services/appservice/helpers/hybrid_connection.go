// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func GetSendKeyValue(ctx context.Context, metadata sdk.ResourceMetaData, id webapps.RelayId, sendKeyName string) (*string, error) {
	relayClient := metadata.Client.Relay.HybridConnectionsClient
 	connectionId := hybridconnections.NewHybridConnectionAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.HybridConnectionNamespaceName, id.RelayName, sendKeyName)
	keys, err := relayClient.ListKeys(ctx, connectionId)
	if err != nil {
		return nil, fmt.Errorf("listing Send Keys for %s in %s: %+v", connectionId, id, err)
	}
	if err != nil || keys.Model == nil || keys.Model.PrimaryKey == nil {
		return nil, fmt.Errorf("reading Send Key Value for %s in %s", connectionId.AuthorizationRuleName, id)
	}
	return keys.Model.PrimaryKey, nil
}
