// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package bot_channel_line

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func botChannelLineSchema() map[string]*pluginsdk.Schema {
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

		"line_channel": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"access_token": {
						Description:  "The access token which is used to call the Line Channel API.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"secret": {
						Description:  "The secret which is used to access the Line Channel.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}
