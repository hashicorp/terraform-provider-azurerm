// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ComponentUpgradeV1ToV2{}

type ComponentUpgradeV1ToV2 struct{}

func (ComponentUpgradeV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return componentSchemaForV1AndV2()
}

func (ComponentUpgradeV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// This state migration is identical to v0 -> v1, however we need to apply it again because application insights
		// resources with the incorrect casing could still be imported and exist within some user's state
		oldIdRaw := rawState["id"].(string)
		id, err := components.ParseComponentIDInsensitively(oldIdRaw)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()

		log.Printf("[DEBUG] Updating ID from %q to %q", oldIdRaw, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}

func componentSchemaForV1AndV2() map[string]*pluginsdk.Schema {
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

		"internet_ingestion_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"internet_query_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"force_customer_storage_for_profiler": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
}
