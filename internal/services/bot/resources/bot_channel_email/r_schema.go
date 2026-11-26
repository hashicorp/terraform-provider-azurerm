// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package bot_channel_email

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func botChannelEmailSchema() map[string]*pluginsdk.Schema {
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

		"email_address": {
			Description:  "The email address that the Bot will authenticate with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"email_password": {
			Description:  "The email password that the Bot will authenticate with.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"email_password", "magic_code"},
		},

		"magic_code": {
			Description:  "The magic code used to set up OAuth authentication.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"email_password", "magic_code"},
		},
	}
}
