package systemcentervirtualmachinemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachineinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SystemCenterVirtualMachineManagerVirtualMachineInstanceModel struct {
	ScopeId                                             string      `tfschema:"scope_id"`
	CustomLocationId                                    string      `tfschema:"custom_location_id"`
	SystemCenterVirtualMachineManagerAvailabilitySetIds []string    `tfschema:"system_center_virtual_machine_manager_availability_set_id"`
	OSProfile                                           []OSProfile `tfschema:"os_profile"`
}

type OSProfile struct {
	ComputerName  string `tfschema:"computer_name"`
	AdminPassword string `tfschema:"admin_password"`
}

var _ sdk.Resource = SystemCenterVirtualMachineManagerVirtualMachineInstanceResource{}
var _ sdk.ResourceWithUpdate = SystemCenterVirtualMachineManagerVirtualMachineInstanceResource{}

type SystemCenterVirtualMachineManagerVirtualMachineInstanceResource struct{}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) ModelObject() interface{} {
	return &SystemCenterVirtualMachineManagerVirtualMachineInstanceModel{}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SystemCenterVirtualMachineManagerVirtualMachineInstanceID
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) ResourceType() string {
	return "azurerm_system_center_virtual_machine_manager_virtual_machine_instance"
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scope_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"custom_location_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"os_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"computer_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"admin_password": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},
				},
			},
		},

		"system_center_virtual_machine_manager_availability_set_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualMachineInstances

			var model SystemCenterVirtualMachineManagerVirtualMachineInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := parse.NewSystemCenterVirtualMachineManagerVirtualMachineInstanceID(model.ScopeId)

			existing, err := client.Get(ctx, id.ScopeId)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := virtualmachineinstances.VirtualMachineInstance{
				ExtendedLocation: virtualmachineinstances.ExtendedLocation{
					Type: utils.String("customLocation"),
					Name: utils.String(model.CustomLocationId),
				},
				Properties: virtualmachineinstances.VirtualMachineInstanceProperties{
					AvailabilitySets: expandVirtualMachineInstanceAvailabilitySets(model.SystemCenterVirtualMachineManagerAvailabilitySetIds),
					OsProfile:        expandVirtualMachineInstanceOSProfile(model.OSProfile),
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id.ScopeId, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualMachineInstances

			id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ScopeId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := SystemCenterVirtualMachineManagerVirtualMachineInstanceModel{}
			if model := resp.Model; model != nil {
				state.ScopeId = id.ScopeId.Scope
				state.CustomLocationId = pointer.From(model.ExtendedLocation.Name)
				state.OSProfile = flattenVirtualMachineInstanceOSProfile(model.Properties.OsProfile)
				state.SystemCenterVirtualMachineManagerAvailabilitySetIds = flattenVirtualMachineInstanceAvailabilitySets(model.Properties.AvailabilitySets)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualMachineInstances

			id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SystemCenterVirtualMachineManagerVirtualMachineInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := &virtualmachineinstances.VirtualMachineInstanceUpdate{
				Properties: &virtualmachineinstances.VirtualMachineInstanceUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("system_center_virtual_machine_manager_availability_set_ids") {
				parameters.Properties.AvailabilitySets = expandVirtualMachineInstanceAvailabilitySets(model.SystemCenterVirtualMachineManagerAvailabilitySetIds)
			}

			if err := client.UpdateThenPoll(ctx, id.ScopeId, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualMachineInstances

			id, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, id.ScopeId, virtualmachineinstances.DeleteOperationOptions{
				DeleteFromHost: pointer.To(virtualmachineinstances.DeleteFromHostTrue),
				Force:          pointer.To(virtualmachineinstances.ForceTrue),
			}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandVirtualMachineInstanceOSProfile(input []OSProfile) *virtualmachineinstances.OsProfileForVMInstance {
	if len(input) == 0 {
		return nil
	}

	osProfile := &input[0]

	result := &virtualmachineinstances.OsProfileForVMInstance{
		ComputerName:  pointer.To(osProfile.ComputerName),
		AdminPassword: pointer.To(osProfile.AdminPassword),
	}

	return result
}

func flattenVirtualMachineInstanceOSProfile(input *virtualmachineinstances.OsProfileForVMInstance) []OSProfile {
	result := make([]OSProfile, 0)
	if input == nil {
		return result
	}

	osProfile := OSProfile{
		ComputerName:  pointer.From(input.ComputerName),
		AdminPassword: pointer.From(input.AdminPassword),
	}

	return append(result, osProfile)
}

func expandVirtualMachineInstanceAvailabilitySets(input []string) *[]virtualmachineinstances.AvailabilitySetListAvailabilitySetsInlined {
	result := make([]virtualmachineinstances.AvailabilitySetListAvailabilitySetsInlined, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		availabilitySet := virtualmachineinstances.AvailabilitySetListAvailabilitySetsInlined{
			Id: pointer.To(v),
		}
		result = append(result, availabilitySet)
	}

	return &result
}

func flattenVirtualMachineInstanceAvailabilitySets(input *[]virtualmachineinstances.AvailabilitySetListAvailabilitySetsInlined) []string {
	result := make([]string, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		if avsId := v.Id; avsId != nil {
			result = append(result, *avsId)
		}
	}

	return result
}
