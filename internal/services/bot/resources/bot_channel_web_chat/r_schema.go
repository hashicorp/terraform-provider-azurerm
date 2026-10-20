// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_web_chat

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func botChannelWebChatSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"bot_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.BotName,
		},

		"site": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"user_upload_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"endpoint_parameters_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"storage_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},
	}
}
