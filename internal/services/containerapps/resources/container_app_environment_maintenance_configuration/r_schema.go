// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_maintenance_configuration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: maintenanceconfigurations.ValidateManagedEnvironmentID,
			Description:  "The ID of the Container App Environment to which this Maintenance Configuration belongs.",
		},

		"maintenance_window": {
			Type:        pluginsdk.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "A `maintenance_window` block as defined below.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"day_of_week": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(maintenanceconfigurations.PossibleValuesForWeekDay(), false),
						Description:  "The day of the week for the maintenance window. Possible values are `Friday`, `Monday`, `Saturday`, `Sunday`, `Thursday`, `Tuesday`, and `Wednesday`.",
					},
					"start_hour_in_utc": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 23),
						Description:  "The start hour of the maintenance window in UTC. Possible values are between `0` and `23`.",
					},
					"duration_hours": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(8, 24),
						Description:  "The duration of the maintenance window in hours. Possible values are between `8` and `24`.",
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
