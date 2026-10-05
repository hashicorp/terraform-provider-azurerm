// Copyright IBM Corp. 2014, 2025
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringDoesNotContainAny("/"),
			),
		},

		"virtual_machine_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"publisher": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"type_handler_version": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"auto_upgrade_minor_version": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"automatic_upgrade_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"failure_suppression_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"settings": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		// due to the sensitive nature, these are not returned by the API
		"protected_settings": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Sensitive:        true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			ConflictsWith:    []string{"protected_settings_from_key_vault"},
		},

		"protected_settings_from_key_vault": helpers.ProtectedSettingsFromKeyVaultSchema(true),

		"provision_after_extensions": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"tags": commonschema.Tags(),
	}
}
