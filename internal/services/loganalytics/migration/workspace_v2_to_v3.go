// Copyright (c) HashiCorp, Inc.
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
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},

		"cmk_for_query_forced": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},

		"daily_quota_gb": {
			Optional: true,
			Type:     pluginsdk.TypeFloat,
		},

		"internet_ingestion_enabled": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},

		"internet_query_enabled": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},

		"local_authentication_disabled": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},

		"location": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"primary_shared_key": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},

		"reservation_capacity_in_gb_per_day": {
			Computed: true,
			Optional: true,
			Type:     pluginsdk.TypeInt,
		},

		"resource_group_name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"retention_in_days": {
			Computed: true,
			Optional: true,
			Type:     pluginsdk.TypeInt,
		},

		"secondary_shared_key": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},

		"sku": {
			Computed: true,
			Optional: true,
			Type:     pluginsdk.TypeString,
		},

		"tags": {
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Optional: true,
			Type:     pluginsdk.TypeMap,
		},

		"workspace_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (WorkspaceV2ToV3) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
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
