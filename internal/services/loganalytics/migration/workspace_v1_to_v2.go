// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = WorkspaceV1ToV2{}

type WorkspaceV1ToV2 struct{}

func (WorkspaceV1ToV2) Schema() map[string]*pluginsdk.Schema {
	// V1 to V2 is the same as v0 to v1 - to workaround a historical issue where `resource_group` was
	// used in place of `resource_group_name` - ergo using the same schema is fine.
	return workspaceSchemaForV0AndV1()
}

func (WorkspaceV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		subscriptionId := meta.(*clients.Client).Account.SubscriptionId

		log.Printf("[DEBUG] Migrating IDs to correct casing for Log Analytics Workspace")
		name := rawState["name"].(string)
		resourceGroup := rawState["resource_group_name"].(string)
		id := workspaces.NewWorkspaceID(subscriptionId, resourceGroup, name)

		rawState["id"] = id.ID()
		return rawState, nil
	}
}
