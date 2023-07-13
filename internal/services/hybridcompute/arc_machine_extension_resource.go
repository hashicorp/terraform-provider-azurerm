// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hybridcompute

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machineextensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MachineExtensionModel struct {
	Name                   string            `tfschema:"name"`
	HybridComputeMachineId string            `tfschema:"arc_machine_id"`
	EnableAutomaticUpgrade bool              `tfschema:"automatic_upgrade_enabled"`
	ForceUpdateTag         string            `tfschema:"force_update_tag"`
	Location               string            `tfschema:"location"`
	ProtectedSettings      string            `tfschema:"protected_settings"`
	Publisher              string            `tfschema:"publisher"`
	Settings               string            `tfschema:"settings"`
	Tags                   map[string]string `tfschema:"tags"`
	Type                   string            `tfschema:"type"`
	TypeHandlerVersion     string            `tfschema:"type_handler_version"`
}

type ArcMachineExtensionResource struct{}

var _ sdk.ResourceWithUpdate = ArcMachineExtensionResource{}

func (r ArcMachineExtensionResource) ResourceType() string {
	return "azurerm_arc_machine_extension"
}

func (r ArcMachineExtensionResource) ModelObject() interface{} {
	return &MachineExtensionModel{}
}

func (r ArcMachineExtensionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return machineextensions.ValidateExtensionID
}

func (r ArcMachineExtensionResource) Arguments() map[string]*pluginsdk.Schema {
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

		"location": commonschema.Location(),

		"arc_machine_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: machines.ValidateMachineID,
		},

		"automatic_upgrade_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"force_update_tag": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"protected_settings": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Sensitive:        true,
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
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"type_handler_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
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

func (r ArcMachineExtensionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ArcMachineExtensionResource) Create() sdk.ResourceFunc {
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

func (r ArcMachineExtensionResource) Update() sdk.ResourceFunc {
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

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ArcMachineExtensionResource) Read() sdk.ResourceFunc {
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

			state := MachineExtensionModel{
				Name:                   id.ExtensionName,
				HybridComputeMachineId: machines.NewMachineID(id.SubscriptionId, id.ResourceGroupName, id.MachineName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				if properties := model.Properties; properties != nil {
					if properties.EnableAutomaticUpgrade != nil {
						state.EnableAutomaticUpgrade = *properties.EnableAutomaticUpgrade
					}

					if properties.ForceUpdateTag != nil {
						state.ForceUpdateTag = *properties.ForceUpdateTag
					}

					if properties.Publisher != nil {
						state.Publisher = *properties.Publisher
					}

					var extModel MachineExtensionModel
					err := metadata.Decode(&extModel)
					if err != nil {
						return err
					}

					if extModel.ProtectedSettings != "" {
						state.ProtectedSettings = extModel.ProtectedSettings
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

			}

			return metadata.Encode(&state)
		},
	}
}

func (r ArcMachineExtensionResource) Delete() sdk.ResourceFunc {
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
