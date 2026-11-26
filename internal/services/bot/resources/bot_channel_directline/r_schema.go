// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package bot_channel_directline

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func botChannelDirectlineSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"bot_name": {
			Description:  "The name of the Bot Resource this channel will be associated with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"site": {
			Description: "A site represents a client application that you want to connect to your bot.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the site.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"user_upload_enabled": {
						Description: "Whether the user upload is enabled for this site.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"enabled": {
						Description: "Enables/Disables this site.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"endpoint_parameters_enabled": {
						Description: "Whether the endpoint parameters are enabled for this site.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"storage_enabled": {
						Description: "Whether the storage site is enabled for detailed logging.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"v1_allowed": {
						Description: "Enables v1 of the Directline protocol for this site.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"v3_allowed": {
						Description: "Enables v3 of the Directline protocol for this site.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"enhanced_authentication_enabled": {
						Description: "Enables additional security measures for this site, see [Direct Line enhanced authentication](https://learn.microsoft.com/azure/bot-service/bot-builder-security-enhanced?view=azure-bot-service-4.0). Disabled by default.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"trusted_origins": {
						Description: "This field is required when `enhanced_authentication_enabled` is enabled. Determines which origins can establish a Directline conversation for this site.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"key": {
						Description: "Primary key for accessing this site.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},

					"key2": {
						Description: "Secondary key for accessing this site.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},

					"id": {
						Description: "ID for the site.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"extension_key_1": {
			Description: "The first key for Directline App Service Extension.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"extension_key_2": {
			Description: "The second key for Directline App Service Extension.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
	}
}
