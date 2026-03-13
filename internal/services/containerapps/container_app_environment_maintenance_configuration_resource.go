// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppEnvironmentMaintenanceConfigurationResource struct{}

type ContainerAppEnvironmentMaintenanceConfigurationModel struct {
	ContainerAppEnvironmentId string              `tfschema:"container_app_environment_id"`
	MaintenanceWindows        []MaintenanceWindow `tfschema:"maintenance_window"`
}

type MaintenanceWindow struct {
	DayOfWeek     string `tfschema:"day_of_week"`
	StartHourInUtc int64  `tfschema:"start_hour_in_utc"`
	DurationHours int64  `tfschema:"duration_hours"`
}

var _ sdk.ResourceWithUpdate = ContainerAppEnvironmentMaintenanceConfigurationResource{}

func (r ContainerAppEnvironmentMaintenanceConfigurationResource) ModelObject() interface{} {
	return &ContainerAppEnvironmentMaintenanceConfigurationModel{}
}

func (r ContainerAppEnvironmentMaintenanceConfigurationResource) ResourceType() string {
	return "azurerm_container_app_environment_maintenance_configuration"
}

func (r ContainerAppEnvironmentMaintenanceConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return maintenanceconfigurations.ValidateMaintenanceConfigurationID
}

func (r ContainerAppEnvironmentMaintenanceConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
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
			MinItems:    1,
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

func (r ContainerAppEnvironmentMaintenanceConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ContainerAppEnvironmentMaintenanceConfigurationResource) Create() sdk.ResourceFunc {
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

func (r ContainerAppEnvironmentMaintenanceConfigurationResource) Read() sdk.ResourceFunc {
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

			var state ContainerAppEnvironmentMaintenanceConfigurationModel

			state.ContainerAppEnvironmentId = maintenanceconfigurations.NewManagedEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName).ID()

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					state.MaintenanceWindows = flattenMaintenanceWindow(props.ScheduledEntries)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ContainerAppEnvironmentMaintenanceConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.MaintenanceConfigurationsClient

			id, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ContainerAppEnvironmentMaintenanceConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.MaintenanceConfigurationsClient

			id, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model or properties was nil", *id)
			}

			maintenanceConfig := *existing.Model

			var model ContainerAppEnvironmentMaintenanceConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("maintenance_window") {
				maintenanceConfig.Properties.ScheduledEntries = expandMaintenanceWindow(model.MaintenanceWindows)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, maintenanceConfig); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

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

func flattenMaintenanceWindow(input []maintenanceconfigurations.ScheduledEntry) []MaintenanceWindow {
	if len(input) == 0 {
		return []MaintenanceWindow{}
	}

	v := input[0]
	return []MaintenanceWindow{
		{
			DayOfWeek:      string(v.WeekDay),
			StartHourInUtc: v.StartHourUtc,
			DurationHours:  v.DurationHours,
		},
	}
}
