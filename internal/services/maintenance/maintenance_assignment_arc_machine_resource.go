// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package maintenance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

//go:generate go run ../../tools/generator-tests resourceidentity -properties "name" -compare-values "subscription_id:arc_machine_id,resource_group_name:arc_machine_id,machine_name:arc_machine_id"

var (
	_ sdk.Resource             = MaintenanceAssignmentArcMachineResource{}
	_ sdk.ResourceWithIdentity = MaintenanceAssignmentArcMachineResource{}
)

type MaintenanceAssignmentArcMachineResource struct{}

type MaintenanceAssignmentArcMachineModel struct {
	ArcMachineId               string `tfschema:"arc_machine_id"`
	MaintenanceConfigurationId string `tfschema:"maintenance_configuration_id"`
	Name                       string `tfschema:"name"`
}

func (MaintenanceAssignmentArcMachineResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"arc_machine_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: machines.ValidateMachineID,
			// The service is returning `arc_machine_id` in lower case, tracked by https://github.com/Azure/azure-rest-api-specs/issues/34824
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"maintenance_configuration_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: maintenanceconfigurations.ValidateMaintenanceConfigurationID,
			// The service is returning `maintenance_configuration_id` in lower case, tracked by https://github.com/Azure/azure-rest-api-specs/issues/34824
			DiffSuppressFunc: suppress.CaseDifference,
		},
	}
}

func (MaintenanceAssignmentArcMachineResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (MaintenanceAssignmentArcMachineResource) ModelObject() any {
	return &MaintenanceAssignmentArcMachineModel{}
}

func (MaintenanceAssignmentArcMachineResource) ResourceType() string {
	return "azurerm_maintenance_assignment_arc_machine"
}

func (MaintenanceAssignmentArcMachineResource) Identity() resourceids.ResourceId {
	return &parse.MaintenanceAssignmentArcMachineId{}
}

func (r MaintenanceAssignmentArcMachineResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Maintenance.ConfigurationAssignmentsClient
			machineClient := metadata.Client.HybridCompute.HybridComputeClient_v2024_07_10.Machines

			var model MaintenanceAssignmentArcMachineModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			maintenanceConfigurationId, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(model.MaintenanceConfigurationId)
			if err != nil {
				return err
			}

			arcMachineId, err := machines.ParseMachineID(model.ArcMachineId)
			if err != nil {
				return err
			}

			id := parse.NewMaintenanceAssignmentArcMachineID(arcMachineId.SubscriptionId, arcMachineId.ResourceGroupName, arcMachineId.MachineName, maintenanceConfigurationId.MaintenanceConfigurationName)
			assignmentId := configurationassignments.NewScopedConfigurationAssignmentID(arcMachineId.ID(), id.ConfigurationAssignmentName)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, assignmentId)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			// The assignment location is required by the API and must match the Arc Machine.
			machine, err := machineClient.Get(ctx, *arcMachineId, machines.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", arcMachineId, err)
			}
			if machine.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", arcMachineId)
			}

			configurationAssignment := configurationassignments.ConfigurationAssignment{
				Name:     pointer.To(maintenanceConfigurationId.MaintenanceConfigurationName),
				Location: pointer.To(location.Normalize(machine.Model.Location)),
				Properties: &configurationassignments.ConfigurationAssignmentProperties{
					MaintenanceConfigurationId: pointer.To(maintenanceConfigurationId.ID()),
					ResourceId:                 pointer.To(arcMachineId.ID()),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, assignmentId, configurationAssignment); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r MaintenanceAssignmentArcMachineResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Maintenance.ConfigurationAssignmentsClient

			id, err := parse.MaintenanceAssignmentArcMachineID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			arcMachineId := machines.NewMachineID(id.SubscriptionId, id.ResourceGroupName, id.MachineName)
			assignmentId := configurationassignments.NewScopedConfigurationAssignmentID(arcMachineId.ID(), id.ConfigurationAssignmentName)
			resp, err := client.Get(ctx, assignmentId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (MaintenanceAssignmentArcMachineResource) flatten(metadata sdk.ResourceMetaData, id *parse.MaintenanceAssignmentArcMachineId, model *configurationassignments.ConfigurationAssignment) error {
	arcMachineId := machines.NewMachineID(id.SubscriptionId, id.ResourceGroupName, id.MachineName)

	state := MaintenanceAssignmentArcMachineModel{
		ArcMachineId: arcMachineId.ID(),
		Name:         id.ConfigurationAssignmentName,
	}

	if props := model.Properties; props != nil && props.MaintenanceConfigurationId != nil {
		configurationId, err := maintenanceconfigurations.ParseMaintenanceConfigurationIDInsensitively(*props.MaintenanceConfigurationId)
		if err != nil {
			return err
		}
		state.MaintenanceConfigurationId = configurationId.ID()
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}

func (r MaintenanceAssignmentArcMachineResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Maintenance.ConfigurationAssignmentsClient

			id, err := parse.MaintenanceAssignmentArcMachineID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			arcMachineId := machines.NewMachineID(id.SubscriptionId, id.ResourceGroupName, id.MachineName)
			assignmentId := configurationassignments.NewScopedConfigurationAssignmentID(arcMachineId.ID(), id.ConfigurationAssignmentName)
			if resp, err := client.Delete(ctx, assignmentId); err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r MaintenanceAssignmentArcMachineResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(input any, key string) (warnings []string, errors []error) {
		v, ok := input.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", key))
			return
		}

		if _, err := parse.MaintenanceAssignmentArcMachineID(v); err != nil {
			errors = append(errors, err)
		}

		return
	}
}
