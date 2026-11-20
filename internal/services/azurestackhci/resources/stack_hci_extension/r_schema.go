// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stack_hci_extension

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/extensions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description: "The name which should be used for this Azure Stack HCI Extension.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringDoesNotContainAny("/"),
			),
		},

		"arc_setting_id": {
			Description:  "The ID of the Azure Stack HCI Cluster Arc Setting.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: extensions.ValidateArcSettingID,
		},

		"publisher": {
			Description:  "The name of the extension handler publisher, such as `Microsoft.Azure.Monitor`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"type": {
			Description:  "Specifies the type of the extension. For example `CustomScriptExtension` or `AzureMonitorLinuxAgent`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"auto_upgrade_minor_version_enabled": {
			Description: "Indicates whether the extension should use a newer minor version if one is available at deployment time. Once deployed, however, the extension will not upgrade minor versions unless redeployed, even with this property set to true.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
			ForceNew:    true,
		},

		"automatic_upgrade_enabled": {
			Description: "Indicates whether the extension should be automatically upgraded by the platform if there is a newer version available.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"protected_settings": {
			Description:      "The json formatted protected settings for the extension.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Sensitive:        true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"settings": {
			Description:      "The json formatted public settings for the extension.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"type_handler_version": {
			Description:  "Specifies the version of the script handler.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config StackHCIExtensionResourceModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if config.AutomaticUpgradeEnabled && config.TypeHandlerVersion != "" {
				return fmt.Errorf("`type_handler_version` cannot be set if `automatic_upgrade_enabled` is true")
			}

			return nil
		},
	}
}
