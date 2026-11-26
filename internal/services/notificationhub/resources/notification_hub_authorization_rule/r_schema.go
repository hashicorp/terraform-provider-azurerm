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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"notification_hub_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"namespace_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"manage": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"send": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"listen": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"primary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}
