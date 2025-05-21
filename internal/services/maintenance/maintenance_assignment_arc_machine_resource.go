package maintenance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

var _ sdk.Resource = MaintenanceAssignmentArcMachineResource{}

type MaintenanceAssignmentArcMachineResource struct{}

type MaintenanceAssignmentArcMachineModel struct {
	Location                   string `tfschema:"location"`
	MaintenanceConfigurationId string `tfschema:"maintenance_configuration_id"`
	ArcMachineId               string `tfschema:"arc_machine_id"`
}

func (MaintenanceAssignmentArcMachineResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

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

func (MaintenanceAssignmentArcMachineResource) ModelObject() interface{} {
	return &MaintenanceAssignmentArcMachineModel{}
}

func (MaintenanceAssignmentArcMachineResource) ResourceType() string {
	return "azurerm_maintenance_assignment_arc_machine"
}

func (r MaintenanceAssignmentArcMachineResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Maintenance.ConfigurationAssignmentsClient

			var model MaintenanceAssignmentArcMachineModel
			if err := metadata.Decode(&model); err != nil {
				return err
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

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			configurationAssignment := configurationassignments.ConfigurationAssignment{
				Name:     &maintenanceConfigurationId.MaintenanceConfigurationName,
				Location: pointer.To(location.Normalize(model.Location)),
				Properties: &configurationassignments.ConfigurationAssignmentProperties{
					MaintenanceConfigurationId: pointer.To(maintenanceConfigurationId.ID()),
					ResourceId:                 pointer.To(arcMachineId.ID()),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, configurationAssignment); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
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

			state := MaintenanceAssignmentArcMachineModel{}

			if model := resp.Model; model != nil {
				// The service is not returning location, tracked by https://github.com/Azure/azure-rest-api-specs/issues/28880
				loc := location.Normalize(pointer.From(model.Location))
				if loc == "" {
					loc = location.Normalize(metadata.ResourceData.Get("location").(string))
				}
				state.Location = loc

				if prop := model.Properties; prop != nil {
					state.ArcMachineId = pointer.From(prop.ResourceId)
					state.MaintenanceConfigurationId = pointer.From(prop.MaintenanceConfigurationId)
				}
			}
			return metadata.Encode(&state)
		},
	}
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

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r MaintenanceAssignmentArcMachineResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(val interface{}, key string) (warns []string, errs []error) {
		idRaw, ok := val.(string)
		if !ok {
			errs = append(errs, fmt.Errorf("expected `id` to be a string but got %+v", val))
			return
		}

		parsedAssignmentId, err := configurationassignments.ParseScopedConfigurationAssignmentID(idRaw)
		if err != nil {
			errs = append(errs, err)
		}

		_, err = machines.ParseMachineIDInsensitively(parsedAssignmentId.Scope)
		if err != nil {
			errs = append(errs, err)
		}

		return
	}
}
