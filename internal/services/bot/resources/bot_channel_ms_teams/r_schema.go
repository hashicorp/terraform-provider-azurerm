// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_ms_teams

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func botChannelMsTeamsSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"bot_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"calling_web_hook": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because this field could not update to empty, so add `Computed: true` to avoid diff
			// issue: https://github.com/Azure/azure-rest-api-specs/issues/9809
			Computed:     true,
			ValidateFunc: validate.BotMSTeamsCallingWebHook(),
		},

		"deployment_environment": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "CommercialDeployment",
			ValidateFunc: validation.StringInSlice([]string{
				"CommercialDeployment",
				"GCCModerateDeployment",
			}, false),
		},

		"calling_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}
