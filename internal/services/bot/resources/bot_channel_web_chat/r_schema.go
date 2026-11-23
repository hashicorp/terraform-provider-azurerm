// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Bot Resource this channel will be associated with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.BotName,
		},

		"site": {
			Description: "A site represents a client application that you want to connect to your bot.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the site.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"user_upload_enabled": {
						Description: "Is the user upload enabled for this site?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"endpoint_parameters_enabled": {
						Description: "Is the endpoint parameters enabled for this site?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"storage_enabled": {
						Description: "Is the storage site enabled for detailed logging?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
				},
			},
		},
	}
}
