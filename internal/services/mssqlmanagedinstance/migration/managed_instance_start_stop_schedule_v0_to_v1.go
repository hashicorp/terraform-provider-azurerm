// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = MsSqlManagedInstanceStartStopScheduleV0ToV1{}

type MsSqlManagedInstanceStartStopScheduleV0ToV1 struct{}

func (MsSqlManagedInstanceStartStopScheduleV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_instance_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"schedule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"start_day": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"start_time": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"stop_day": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"stop_time": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"timezone_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "UTC",
		},

		"next_execution_time": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"next_run_action": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (MsSqlManagedInstanceStartStopScheduleV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		return rawState, nil
	}
}
