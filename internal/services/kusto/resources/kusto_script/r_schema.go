// Copyright IBM Corp. 2014, 2025
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"database_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKustoDatabaseID,
		},

		"continue_on_errors_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"force_an_update_when_value_changed": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"url", "script_content"},
			RequiredWith: []string{"sas_token"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sas_token": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"url"},
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"script_content": {
			Type:         pluginsdk.TypeString,
			ExactlyOneOf: []string{"url", "script_content"},
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"script_level": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(scripts.ScriptLevelDatabase),
			ValidateFunc: validation.StringInSlice(scripts.PossibleValuesForScriptLevel(), false),
		},

		"principal_permissions_action": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(scripts.PrincipalPermissionsActionRetainPermissionOnScriptCompletion),
			ValidateFunc: validation.StringInSlice(scripts.PossibleValuesForPrincipalPermissionsAction(), false),
		},
	}
}
