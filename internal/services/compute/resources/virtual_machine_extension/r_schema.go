// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_extension

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualMachineExtensionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Virtual Machine Extension.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringDoesNotContainAny("/"),
			),
		},

		"virtual_machine_id": {
			Description:  "The ID of the Virtual Machine.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"publisher": {
			Description: "The publisher of the extension, available publishers can be found by using the Azure CLI.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"type": {
			Description: "The type of extension, available types for a publisher can be found using the Azure CLI.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},

		"type_handler_version": {
			Description: "The version of the extension to use, available versions can be found using the Azure CLI.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},

		"auto_upgrade_minor_version": {
			Description: "Whether the platform deploys the latest minor version update to the `type_handler_version` specified.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
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

		"settings": {
			Description:      "The settings passed to the extension, these are specified as a JSON object in a string.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		// due to the sensitive nature, these are not returned by the API
		"protected_settings": {
			Description:      "The protected_settings passed to the extension, like settings, these are specified as a JSON object in a string.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Sensitive:        true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			ConflictsWith:    []string{"protected_settings_from_key_vault"},
		},

		"protected_settings_from_key_vault": helpers.ProtectedSettingsFromKeyVaultSchema(true),

		"provision_after_extensions": {
			Description: "The collection of extension names after which this extension needs to be provisioned.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"tags": commonschema.Tags(),
	}
}
