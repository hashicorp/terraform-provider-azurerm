package azurestackhci

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-08-01/extensions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = StackHCIExtensionResource{}
var _ sdk.ResourceWithUpdate = StackHCIExtensionResource{}

type StackHCIExtensionResource struct{}

func (r StackHCIExtensionResource) ResourceType() string {
	return "azurerm_stack_hci_extension"
}

func (r StackHCIExtensionResource) ModelObject() interface{} {
	return &StackHCIExtensionResourceModel{}
}

func (r StackHCIExtensionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return extensions.ValidateExtensionID
}

type StackHCIExtensionResourceModel struct {
	Name                    string `tfschema:"name"`
	ArcSettingId            string `tfschema:"arc_setting_id"`
	AutoUpgradeMinorVersion bool   `tfschema:"auto_upgrade_minor_version"`
	AutomaticUpgradeEnabled bool   `tfschema:"automatic_upgrade_enabled"`
	ForceUpdateTag          string `tfschema:"force_update_tag"`
	ProtectedSettings       string `tfschema:"protected_setting"`
	Publisher               string `tfschema:"publisher"`
	Setting                 string `tfschema:"setting"`
	Type                    string `tfschema:"type"`
	TypeHandlerVersion      string `tfschema:"type_handler_version"`
}

func (r StackHCIExtensionResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringDoesNotContainAny("/"),
			),
		},

		"arc_setting_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: extensions.ValidateArcSettingID,
		},

		"publisher": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"auto_upgrade_minor_version": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"automatic_upgrade_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"force_update_tag": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"protected_setting": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Sensitive:        true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"setting": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"type_handler_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r StackHCIExtensionResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r StackHCIExtensionResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config StackHCIExtensionResourceModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if config.AutomaticUpgradeEnabled == true && config.TypeHandlerVersion != "" {
				return fmt.Errorf("`type_handler_version` cannot be set if `automatic_upgrade_enabled` is true")
			}

			return nil
		},
	}
}

func (r StackHCIExtensionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.Extensions

			var config StackHCIExtensionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			arcSettingId, err := extensions.ParseArcSettingID(config.ArcSettingId)
			if err != nil {
				return err
			}

			id := extensions.NewExtensionID(arcSettingId.SubscriptionId, arcSettingId.ResourceGroupName, arcSettingId.ClusterName, arcSettingId.ArcSettingName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			input := extensions.Extension{
				Properties: &extensions.ExtensionProperties{
					ExtensionParameters: &extensions.ExtensionParameters{
						AutoUpgradeMinorVersion: pointer.To(config.AutoUpgradeMinorVersion),
						EnableAutomaticUpgrade:  pointer.To(config.AutomaticUpgradeEnabled),
						ForceUpdateTag:          pointer.To(config.ForceUpdateTag),
						Publisher:               pointer.To(config.Publisher),
						Type:                    pointer.To(config.Type),
					},
				},
			}

			if config.TypeHandlerVersion != "" {
				input.Properties.ExtensionParameters.TypeHandlerVersion = pointer.To(config.TypeHandlerVersion)
			}

			if config.Setting != "" {
				expandedSetting, err := pluginsdk.ExpandJsonFromString(config.Setting)
				if err != nil {
					return fmt.Errorf("expanding `setting`: %+v", err)
				}

				input.Properties.ExtensionParameters.Settings = pointer.To(interface{}(expandedSetting))
			}

			if config.ProtectedSettings != "" {
				expandedSetting, err := pluginsdk.ExpandJsonFromString(config.ProtectedSettings)
				if err != nil {
					return fmt.Errorf("expanding `protected_setting`: %+v", err)
				}

				input.Properties.ExtensionParameters.ProtectedSettings = pointer.To(interface{}(expandedSetting))
			}

			if err := client.CreateThenPoll(ctx, id, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil

		},
	}
}

func (r StackHCIExtensionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.Extensions

			id, err := extensions.ParseExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			// protected_settings is not returned in the response, so we read it from the state
			var extension StackHCIExtensionResourceModel
			if err := metadata.Decode(&extension); err != nil {
				return err
			}

			if model := existing.Model; model != nil {
				extension.Name = id.ExtensionName
				extension.ArcSettingId = extensions.NewArcSettingID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ArcSettingName).ID()

				if props := model.Properties; props != nil {
					if param := props.ExtensionParameters; param != nil {
						extension.AutomaticUpgradeEnabled = pointer.From(param.EnableAutomaticUpgrade)
						extension.AutoUpgradeMinorVersion = pointer.From(param.AutoUpgradeMinorVersion)
						extension.ForceUpdateTag = pointer.From(param.ForceUpdateTag)
						extension.Publisher = pointer.From(param.Publisher)
						extension.Type = pointer.From(param.Type)
						extension.TypeHandlerVersion = pointer.From(param.TypeHandlerVersion)

						var setting string
						if param.Settings != nil {
							setting, err = pluginsdk.FlattenJsonToString((*param.Settings).(map[string]interface{}))
							if err != nil {
								return fmt.Errorf("flatenning `settings`: %+v", err)
							}
						}
						extension.Setting = setting
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&extension)
		},
	}
}

func (r StackHCIExtensionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.Extensions

			id, err := extensions.ParseExtensionID(metadata.ResourceData.Id())
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

func (r StackHCIExtensionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.Extensions

			var config StackHCIExtensionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			id, err := extensions.ParseExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil || model.Properties == nil || model.Properties.ExtensionParameters == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("auto_upgrade_enabled") {
				model.Properties.ExtensionParameters.EnableAutomaticUpgrade = pointer.To(config.AutomaticUpgradeEnabled)
			}

			if metadata.ResourceData.HasChange("protected_setting") {
				expandedSetting, err := pluginsdk.ExpandJsonFromString(config.ProtectedSettings)
				if err != nil {
					return fmt.Errorf("expanding `protected_setting`: %+v", err)
				}

				model.Properties.ExtensionParameters.ProtectedSettings = pointer.To(interface{}(expandedSetting))
			}

			if metadata.ResourceData.HasChange("setting") {
				expandedSetting, err := pluginsdk.ExpandJsonFromString(config.Setting)
				if err != nil {
					return fmt.Errorf("expanding `setting`: %+v", err)
				}

				model.Properties.ExtensionParameters.Settings = pointer.To(interface{}(expandedSetting))
			}

			if metadata.ResourceData.HasChange("type_handler_version") {
				model.Properties.ExtensionParameters.TypeHandlerVersion = pointer.To(config.TypeHandlerVersion)
			}

			if metadata.ResourceData.HasChange("auto_upgrade_minor_version") {
				model.Properties.ExtensionParameters.AutoUpgradeMinorVersion = pointer.To(config.AutoUpgradeMinorVersion)
			}

			if err := client.CreateThenPoll(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
