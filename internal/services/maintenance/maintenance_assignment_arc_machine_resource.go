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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

//go:generate go run ../../tools/generator-tests resourceidentity -properties "scope:arc_machine_id" -compare-values "name:maintenance_configuration_id" -has-id-casing-bug

var (
	_ sdk.Resource                   = MaintenanceAssignmentArcMachineResource{}
	_ sdk.ResourceWithIdentity       = MaintenanceAssignmentArcMachineResource{}
	_ sdk.ResourceWithCustomImporter = MaintenanceAssignmentArcMachineResource{}
)

type MaintenanceAssignmentArcMachineResource struct{}

type MaintenanceAssignmentArcMachineModel struct {
	ArcMachineId               string `tfschema:"arc_machine_id"`
	MaintenanceConfigurationId string `tfschema:"maintenance_configuration_id"`
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
	return map[string]*pluginsdk.Schema{}
}

func (MaintenanceAssignmentArcMachineResource) ModelObject() any {
	return &MaintenanceAssignmentArcMachineModel{}
}

func (MaintenanceAssignmentArcMachineResource) ResourceType() string {
	return "azurerm_maintenance_assignment_arc_machine"
}

func (MaintenanceAssignmentArcMachineResource) Identity() resourceids.ResourceId {
	return &configurationassignments.ScopedConfigurationAssignmentId{}
}

func (MaintenanceAssignmentArcMachineResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := configurationassignments.ParseScopedConfigurationAssignmentID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		// The generic scoped ID also accepts assignments on VMs and Dedicated Hosts.
		// Validate the parent so both ID and identity imports only accept Arc Machine assignments.
		if _, err := machines.ParseMachineID(id.Scope); err != nil {
			return err
		}

		return nil
	}
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

			id := configurationassignments.NewScopedConfigurationAssignmentID(arcMachineId.ID(), maintenanceConfigurationId.MaintenanceConfigurationName)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
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

			if _, err := client.CreateOrUpdate(ctx, id, configurationAssignment); err != nil {
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

			id, err := configurationassignments.ParseScopedConfigurationAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
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

func (MaintenanceAssignmentArcMachineResource) flatten(metadata sdk.ResourceMetaData, id *configurationassignments.ScopedConfigurationAssignmentId, model *configurationassignments.ConfigurationAssignment) error {
	arcMachineId, err := machines.ParseMachineIDInsensitively(id.Scope)
	if err != nil {
		return err
	}
	id.Scope = arcMachineId.ID()

	state := MaintenanceAssignmentArcMachineModel{
		ArcMachineId: arcMachineId.ID(),
	}

	if props := model.Properties; props != nil && props.MaintenanceConfigurationId != nil {
		configurationId, err := maintenanceconfigurations.ParseMaintenanceConfigurationIDInsensitively(*props.MaintenanceConfigurationId)
		if err != nil {
			return err
		}
		state.MaintenanceConfigurationId = configurationId.ID()
	}

	metadata.SetID(id)
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

			id, err := configurationassignments.ParseScopedConfigurationAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if resp, err := client.Delete(ctx, *id); err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r MaintenanceAssignmentArcMachineResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return configurationassignments.ValidateScopedConfigurationAssignmentID
}
