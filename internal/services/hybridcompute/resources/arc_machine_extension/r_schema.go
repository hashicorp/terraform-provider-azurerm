// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_machine_extension

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name which should be used for this Hybrid Compute Machine Extension.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringDoesNotContainAny("/"),
			),
		},

		"location": commonschema.Location(),

		"arc_machine_id": {
			Description:  "The ID of the Hybrid Compute Machine Extension.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: machines.ValidateMachineID,
		},

		"automatic_upgrade_enabled": {
			Description: "Whether the extension should be automatically upgraded by the platform if there is a newer version available.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"force_update_tag": {
			Description:  "How the extension handler should be forced to update even if the extension configuration has not changed.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"protected_settings": {
			Description:      "JSON formatted protected settings for the extension.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Sensitive:        true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"publisher": {
			Description:  "The name of the extension handler publisher, such as `Microsoft.Azure.Monitor`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"settings": {
			Description:      "JSON formatted public settings for the extension.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"type": {
			Description:  "The type of the extension. For example `CustomScriptExtension` or `AzureMonitorLinuxAgent`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"type_handler_version": {
			Description:  "The version of the script handler.",
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

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
