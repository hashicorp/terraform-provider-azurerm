package hybridcompute

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-03-10/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machineextensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MachineExtensionModel struct {
	Name                   string                              `tfschema:"name"`
	HybridComputeMachineId string                              `tfschema:"hybrid_compute_machine_id"`
	EnableAutomaticUpgrade bool                                `tfschema:"automatic_upgrade_enabled"`
	ForceUpdateTag         string                              `tfschema:"force_update_tag"`
	InstanceView           []MachineExtensionInstanceViewModel `tfschema:"instance_view"`
	Location               string                              `tfschema:"location"`
	ProtectedSettings      string                              `tfschema:"protected_settings"`
	Publisher              string                              `tfschema:"publisher"`
	Settings               string                              `tfschema:"settings"`
	Tags                   map[string]string                   `tfschema:"tags"`
	Type                   string                              `tfschema:"type"`
	TypeHandlerVersion     string                              `tfschema:"type_handler_version"`
}

type MachineExtensionInstanceViewModel struct {
	Name               string                                    `tfschema:"name"`
	Status             []MachineExtensionInstanceViewStatusModel `tfschema:"status"`
	Type               string                                    `tfschema:"type"`
	TypeHandlerVersion string                                    `tfschema:"type_handler_version"`
}

type MachineExtensionInstanceViewStatusModel struct {
	Code          string                             `tfschema:"code"`
	DisplayStatus string                             `tfschema:"display_status"`
	Level         machineextensions.StatusLevelTypes `tfschema:"level"`
	Message       string                             `tfschema:"message"`
	Time          string                             `tfschema:"time"`
}

type HybridComputeMachineExtensionResource struct{}

var _ sdk.ResourceWithUpdate = HybridComputeMachineExtensionResource{}

func (r HybridComputeMachineExtensionResource) ResourceType() string {
	return "azurerm_hybrid_compute_machine_extension"
}

func (r HybridComputeMachineExtensionResource) ModelObject() interface{} {
	return &MachineExtensionModel{}
}

func (r HybridComputeMachineExtensionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return machineextensions.ValidateExtensionID
}

func (r HybridComputeMachineExtensionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringDoesNotContainAny("/"),
			),
		},

		"hybrid_compute_machine_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: machines.ValidateMachineID,
		},

		"automatic_upgrade_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"force_update_tag": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"protected_settings": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"publisher": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"settings": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"type_handler_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				// suppress any diff if automatic_upgrade_enabled is true
				if value, ok := d.GetOk("automatic_upgrade_enabled"); ok && value.(bool) {
					return true
				}
				// e.g. 1.24 -> 1.24.1 will be considered as no change
				if len(oldValue) > 0 && len(newValue) > 0 && strings.HasPrefix(oldValue, newValue) {
					return true
				}
				return false
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r HybridComputeMachineExtensionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"instance_view": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"status": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"code": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"display_status": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"level": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"message": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"time": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"type_handler_version": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r HybridComputeMachineExtensionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MachineExtensionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.HybridCompute.MachineExtensionsClient
			machineId, err := machines.ParseMachineID(model.HybridComputeMachineId)
			if err != nil {
				return err
			}

			id := machineextensions.NewExtensionID(machineId.SubscriptionId, machineId.ResourceGroupName, machineId.MachineName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &machineextensions.MachineExtension{
				Location: location.Normalize(model.Location),
				Properties: &machineextensions.MachineExtensionProperties{
					EnableAutomaticUpgrade: &model.EnableAutomaticUpgrade,
				},
				Tags: &model.Tags,
			}

			if model.ForceUpdateTag != "" {
				properties.Properties.ForceUpdateTag = &model.ForceUpdateTag
			}

			if model.ProtectedSettings != "" {
				var protectedSettingsValue interface{}
				err = json.Unmarshal([]byte(model.ProtectedSettings), &protectedSettingsValue)
				if err != nil {
					return err
				}
				properties.Properties.ProtectedSettings = &protectedSettingsValue
			}

			if model.Publisher != "" {
				properties.Properties.Publisher = &model.Publisher
			}

			if model.Settings != "" {
				var settingsValue interface{}
				err = json.Unmarshal([]byte(model.Settings), &settingsValue)
				if err != nil {
					return err
				}
				properties.Properties.Settings = &settingsValue
			}

			if model.Type != "" {
				properties.Properties.Type = &model.Type
			}

			if model.TypeHandlerVersion != "" {
				properties.Properties.TypeHandlerVersion = &model.TypeHandlerVersion
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r HybridComputeMachineExtensionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HybridCompute.MachineExtensionsClient

			id, err := machineextensions.ParseExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model MachineExtensionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("automatic_upgrade_enabled") {
				properties.Properties.EnableAutomaticUpgrade = &model.EnableAutomaticUpgrade
			}

			if metadata.ResourceData.HasChange("force_update_tag") {
				if model.ForceUpdateTag != "" {
					properties.Properties.ForceUpdateTag = &model.ForceUpdateTag
				} else {
					properties.Properties.ForceUpdateTag = nil
				}
			}

			if metadata.ResourceData.HasChange("protected_settings") {
				var protectedSettingsValue interface{}
				err := json.Unmarshal([]byte(model.ProtectedSettings), &protectedSettingsValue)
				if err != nil {
					return err
				}

				properties.Properties.ProtectedSettings = &protectedSettingsValue
			}

			if metadata.ResourceData.HasChange("publisher") {
				if model.Publisher != "" {
					properties.Properties.Publisher = &model.Publisher
				} else {
					properties.Properties.Publisher = nil
				}
			}

			if metadata.ResourceData.HasChange("settings") {
				var settingsValue interface{}
				err := json.Unmarshal([]byte(model.Settings), &settingsValue)
				if err != nil {
					return err
				}

				properties.Properties.Settings = &settingsValue
			}

			if metadata.ResourceData.HasChange("type") {
				if model.Type != "" {
					properties.Properties.Type = &model.Type
				} else {
					properties.Properties.Type = nil
				}
			}

			if metadata.ResourceData.HasChange("type_handler_version") {
				if model.TypeHandlerVersion != "" {
					properties.Properties.TypeHandlerVersion = &model.TypeHandlerVersion
				} else {
					properties.Properties.TypeHandlerVersion = nil
				}
			}

			properties.SystemData = nil

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r HybridComputeMachineExtensionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HybridCompute.MachineExtensionsClient

			id, err := machineextensions.ParseExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := MachineExtensionModel{
				Name:                   id.ExtensionName,
				HybridComputeMachineId: machines.NewMachineID(id.SubscriptionId, id.ResourceGroupName, id.MachineName).ID(),
				Location:               location.Normalize(model.Location),
			}

			if properties := model.Properties; properties != nil {
				if properties.EnableAutomaticUpgrade != nil {
					state.EnableAutomaticUpgrade = *properties.EnableAutomaticUpgrade
				}

				if properties.ForceUpdateTag != nil {
					state.ForceUpdateTag = *properties.ForceUpdateTag
				}

				instanceViewValue, err := flattenMachineExtensionInstanceViewModel(properties.InstanceView)
				if err != nil {
					return err
				}

				state.InstanceView = instanceViewValue

				if properties.ProtectedSettings != nil && *properties.ProtectedSettings != nil {

					protectedSettingsValue, err := json.Marshal(*properties.ProtectedSettings)
					if err != nil {
						return err
					}

					state.ProtectedSettings = string(protectedSettingsValue)
				}

				if properties.Publisher != nil {
					state.Publisher = *properties.Publisher
				}

				if properties.Settings != nil && *properties.Settings != nil {

					settingsValue, err := json.Marshal(*properties.Settings)
					if err != nil {
						return err
					}

					state.Settings = string(settingsValue)
				}

				if properties.Type != nil {
					state.Type = *properties.Type
				}

				if properties.TypeHandlerVersion != nil {
					state.TypeHandlerVersion = *properties.TypeHandlerVersion
				}
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r HybridComputeMachineExtensionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HybridCompute.MachineExtensionsClient

			id, err := machineextensions.ParseExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func flattenMachineExtensionInstanceViewModel(input *machineextensions.MachineExtensionInstanceView) ([]MachineExtensionInstanceViewModel, error) {
	var outputList []MachineExtensionInstanceViewModel
	if input == nil {
		return outputList, nil
	}

	output := MachineExtensionInstanceViewModel{}

	if input.Name != nil {
		output.Name = *input.Name
	}

	statusValue, err := flattenMachineExtensionInstanceViewStatusModel(input.Status)
	if err != nil {
		return nil, err
	}

	output.Status = statusValue

	if input.Type != nil {
		output.Type = *input.Type
	}

	if input.TypeHandlerVersion != nil {
		output.TypeHandlerVersion = *input.TypeHandlerVersion
	}

	return append(outputList, output), nil
}

func flattenMachineExtensionInstanceViewStatusModel(input *machineextensions.MachineExtensionInstanceViewStatus) ([]MachineExtensionInstanceViewStatusModel, error) {
	var outputList []MachineExtensionInstanceViewStatusModel
	if input == nil {
		return outputList, nil
	}

	output := MachineExtensionInstanceViewStatusModel{}

	if input.Code != nil {
		output.Code = *input.Code
	}

	if input.DisplayStatus != nil {
		output.DisplayStatus = *input.DisplayStatus
	}

	if input.Level != nil {
		output.Level = *input.Level
	}

	if input.Message != nil {
		output.Message = *input.Message
	}

	if input.Time != nil {
		output.Time = *input.Time
	}

	return append(outputList, output), nil
}
