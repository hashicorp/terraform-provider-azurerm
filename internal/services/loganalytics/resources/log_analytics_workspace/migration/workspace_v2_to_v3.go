// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = WorkspaceV2ToV3{}

type WorkspaceV2ToV3 struct{}

func (WorkspaceV2ToV3) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"allow_resource_only_permissions": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"cmk_for_query_forced": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"daily_quota_gb": {
			Type:     pluginsdk.TypeFloat,
			Optional: true,
		},

		"internet_ingestion_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"internet_query_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"local_authentication_disabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"primary_shared_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"reservation_capacity_in_gb_per_day": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"retention_in_days": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"secondary_shared_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},

		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (WorkspaceV2ToV3) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
		oldId := rawState["id"].(string)
		newId, err := workspaces.ParseWorkspaceIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
