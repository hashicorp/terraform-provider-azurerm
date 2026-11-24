// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_firewall_rule

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mysqlFlexibleServerFirewallRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the MySQL Firewall Rule.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"server_name": {
			Description:  "The name of the MySQL Flexible Server.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FlexibleServerName,
		},

		"start_ip_address": {
			Description:  "The Start IP Address associated with this Firewall Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsIPv4Address,
		},

		"end_ip_address": {
			Description:  "The End IP Address associated with this Firewall Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsIPv4Address,
		},
	}
}
