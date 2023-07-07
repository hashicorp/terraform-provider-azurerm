// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/schedules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = DevTestLabScheduleUpgradeV0ToV1{}

type DevTestLabScheduleUpgradeV0ToV1 struct{}

func (DevTestLabScheduleUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return devTestLabScheduleSchemaForV0AndV1()
}

func (DevTestLabScheduleUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old:
		// 	/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/microsoft.devtestlab/labs/{labName}/schedules/{scheduleName}
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}/schedules/{scheduleName}
		oldId := rawState["id"].(string)
		id, err := schedules.ParseLabScheduleIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}

func devTestLabScheduleSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": commonschema.Location(),

		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"lab_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"task_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"weekly_recurrence": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"time": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"week_days": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"daily_recurrence": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"time": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"hourly_recurrence": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"minute": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
		},

		"time_zone_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"notification_settings": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"status": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"time_in_minutes": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
					"webhook_url": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"tags": tags.Schema(),
	}
}
