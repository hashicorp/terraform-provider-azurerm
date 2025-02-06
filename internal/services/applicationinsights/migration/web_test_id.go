// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = WebTestUpgradeV0ToV1{}

type WebTestUpgradeV0ToV1 struct{}

func (WebTestUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return webTestSchemaForV0AndV1()
}

func (WebTestUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.insights/webtests/test1
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/webTests/test1
		oldIdRaw := rawState["id"].(string)
		id, err := webtestsapis.ParseWebTestIDInsensitively(oldIdRaw)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldIdRaw, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}

func webTestSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"application_insights_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": commonschema.Location(),

		"kind": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"frequency": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"timeout": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"retry_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"geo_locations": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"configuration": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"tags": tags.Schema(),

		"synthetic_monitor_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
