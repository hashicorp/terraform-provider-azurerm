// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_maintenance_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name container_app_environment_maintenance_configuration -service-package-name containerapps -compare-values "subscription_id:container_app_environment_id,resource_group_name:container_app_environment_id,managed_environment_name:container_app_environment_id" -known-values "name:default"

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.MaintenanceConfigurationsClient

			id, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			return r.flatten(metadata, id, existing.Model)
		},
	}
}

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *maintenanceconfigurations.MaintenanceConfigurationId, model *maintenanceconfigurations.MaintenanceConfigurationResource) error {
	state := ContainerAppEnvironmentMaintenanceConfigurationModel{
		ContainerAppEnvironmentId: maintenanceconfigurations.NewManagedEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName).ID(),
	}

	if model != nil {
		if props := model.Properties; props != nil {
			state.MaintenanceWindows = flattenMaintenanceWindow(props.ScheduledEntries)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenMaintenanceWindow(input []maintenanceconfigurations.ScheduledEntry) []MaintenanceWindow {
	if len(input) == 0 {
		return []MaintenanceWindow{}
	}

	v := input[0]
	return []MaintenanceWindow{
		{
			DayOfWeek:       string(v.WeekDay),
			StartHourInUtc:  v.StartHourUtc,
			DurationInHours: v.DurationHours,
		},
	}
}
