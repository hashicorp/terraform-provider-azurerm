// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ComponentUpgradeV0ToV1{}

type ComponentUpgradeV0ToV1 struct{}

func (ComponentUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return componentSchemaForV0AndV1()
}

func (ComponentUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.insights/components/component1
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/component1
		oldIdRaw := rawState["id"].(string)
		id, err := parse.ComponentIDInsensitively(oldIdRaw)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()

		log.Printf("[DEBUG] Updating ID from %q to %q", oldIdRaw, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}

func componentSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"application_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"retention_in_days": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"sampling_percentage": {
			Type:     pluginsdk.TypeFloat,
			Optional: true,
		},

		"disable_ip_masking": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"tags": tags.Schema(),

		"daily_data_cap_in_gb": {
			Type:     pluginsdk.TypeFloat,
			Optional: true,
			Computed: true,
		},

		"daily_data_cap_notifications_disabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"app_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"instrumentation_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"local_authentication_disabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
}
