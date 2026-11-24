// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package notification_hub_authorization_rule

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func notificationHubAuthorizationRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name to use for this Authorization Rule.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"notification_hub_name": {
			Description: "The name of the Notification Hub for which the Authorization Rule should be created.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"namespace_name": {
			Description: "The name of the Notification Hub Namespace in which the Notification Hub exists.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"manage": {
			Description: "Does this Authorization Rule have Manage access to the Notification Hub?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"send": {
			Description: "Does this Authorization Rule have Send access to the Notification Hub?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"listen": {
			Description: "Does this Authorization Rule have Listen access to the Notification Hub?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"primary_access_key": {
			Description: "The Primary Access Key associated with this Authorization Rule.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_access_key": {
			Description: "The Secondary Access Key associated with this Authorization Rule.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"primary_connection_string": {
			Description: "The Primary Connection String associated with this Authorization Rule.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_connection_string": {
			Description: "The Secondary Connection String associated with this Authorization Rule.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
	}
}
