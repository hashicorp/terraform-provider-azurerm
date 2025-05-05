// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/virtualendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PostgresqlFlexibleServerVirtualEndpointV0toV1 struct{}

var _ pluginsdk.StateUpgrade = PostgresqlFlexibleServerVirtualEndpointV0toV1{}

func (l PostgresqlFlexibleServerVirtualEndpointV0toV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"source_server_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"replica_server_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (l PostgresqlFlexibleServerVirtualEndpointV0toV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		name := rawState["name"].(string)

		sourceServerId, err := servers.ParseFlexibleServerID(rawState["source_server_id"].(string))
		if err != nil {
			return nil, err
		}

		sourceEndpointId := virtualendpoints.NewVirtualEndpointID(sourceServerId.SubscriptionId, sourceServerId.ResourceGroupName, sourceServerId.FlexibleServerName, name)

		replicaServerId, err := servers.ParseFlexibleServerID(rawState["replica_server_id"].(string))
		if err != nil {
			return nil, err
		}

		replicaEndpointId := virtualendpoints.NewVirtualEndpointID(replicaServerId.SubscriptionId, replicaServerId.ResourceGroupName, replicaServerId.FlexibleServerName, name)

		newId := commonids.NewCompositeResourceID(&sourceEndpointId, &replicaEndpointId).ID()

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId
		return rawState, nil
	}
}
