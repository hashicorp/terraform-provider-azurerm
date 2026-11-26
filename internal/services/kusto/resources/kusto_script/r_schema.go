// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto_script

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/scripts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func kustoScriptSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name which should be used for this Kusto Script.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"database_id": {
			Description:  "The ID of the Kusto Database.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKustoDatabaseID,
		},

		"continue_on_errors_enabled": {
			Description: "Flag that indicates whether to continue if one of the command fails.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"force_an_update_when_value_changed": {
			Description:  "A unique string. If changed the script will be applied again.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"url": {
			Description:  "The URL to the KQL script blob file. Must not be used together with scriptContent property. Reference [this documentation](https://learn.microsoft.com/azure/data-explorer/database-script) that describes the commands that are allowed in the script.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"sas_token"},
			ExactlyOneOf: []string{"url", "script_content"},
		},

		"sas_token": {
			Description:  "The SAS token used to access the script.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"url"},
		},

		"script_content": {
			Description:  "The script content. This property should be used when the script is provided inline and not through file in an SA. Must not be used together with `url` and `sas_token` properties.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"url", "script_content"},
		},

		"script_level": {
			Description:  "The type of script commands.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(scripts.ScriptLevelDatabase),
			ValidateFunc: validation.StringInSlice(scripts.PossibleValuesForScriptLevel(), false),
		},

		"principal_permissions_action": {
			Description:  "Whether the script caller's permissions remain in effect after the script has finished running.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(scripts.PrincipalPermissionsActionRetainPermissionOnScriptCompletion),
			ValidateFunc: validation.StringInSlice(scripts.PossibleValuesForPrincipalPermissionsAction(), false),
		},
	}
}
