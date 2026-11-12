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
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.MaintenanceConfigurationsClient

			var model ContainerAppEnvironmentMaintenanceConfigurationModel

			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			containerAppEnvironmentId, err := maintenanceconfigurations.ParseManagedEnvironmentID(model.ContainerAppEnvironmentId)
			if err != nil {
				return err
			}

			id := maintenanceconfigurations.NewMaintenanceConfigurationID(metadata.Client.Account.SubscriptionId, containerAppEnvironmentId.ResourceGroupName, containerAppEnvironmentId.ManagedEnvironmentName, "default")

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			maintenanceConfig := maintenanceconfigurations.MaintenanceConfigurationResource{
				Properties: &maintenanceconfigurations.ScheduledEntries{
					ScheduledEntries: expandMaintenanceWindow(model.MaintenanceWindows),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, maintenanceConfig); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func expandMaintenanceWindow(input []MaintenanceWindow) []maintenanceconfigurations.ScheduledEntry {
	if len(input) == 0 {
		return []maintenanceconfigurations.ScheduledEntry{}
	}

	v := input[0]
	return []maintenanceconfigurations.ScheduledEntry{
		{
			WeekDay:       maintenanceconfigurations.WeekDay(v.DayOfWeek),
			StartHourUtc:  v.StartHourInUtc,
			DurationHours: v.DurationHours,
		},
	}
}
