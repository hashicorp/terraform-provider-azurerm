// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set_extension

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualMachineScaleSetExtensionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name for the Virtual Machine Scale Set Extension.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringDoesNotContainAny("/"),
			),
		},

		"virtual_machine_scale_set_id": {
			Description:  "The ID of the Virtual Machine Scale Set.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineScaleSetID,
		},

		"publisher": {
			Description:  "The Publisher of the Extension.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"type": {
			Description:  "The Type of the Extension.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"type_handler_version": {
			Description:  "The version of the extension to use, available versions can be found using the Azure CLI.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"auto_upgrade_minor_version": {
			Description: "Should the latest version of the Extension be used at Deployment Time, if one is available? This won't auto-update the extension on existing installation.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"automatic_upgrade_enabled": {
			Description: "Should the Extension be automatically updated whenever the Publisher releases a new version of this VM Extension?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"failure_suppression_enabled": {
			Description: "Should failures from the extension be suppressed?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"force_update_tag": {
			Description: "A value which, when different to the previous value can be used to force-run the Extension even if the Extension Configuration hasn't changed.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"protected_settings": {
			Description:      "A JSON String which specifies Sensitive Settings (such as Passwords) for the Extension.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Sensitive:        true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"protected_settings_from_key_vault": helpers.ProtectedSettingsFromKeyVaultSchema(true),

		"provision_after_extensions": {
			Description: "An ordered list of Extension names which this should be provisioned after.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"settings": {
			Description:      "A JSON String which specifies Settings for the Extension.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},
	}
}
