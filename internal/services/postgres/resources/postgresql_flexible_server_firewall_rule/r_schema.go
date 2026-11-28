// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_firewall_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2026-10-01/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func postgresqlFlexibleServerFirewallRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this PostgreSQL Flexible Server Firewall Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FlexibleServerFirewallRuleName,
		},

		"server_id": {
			Description:  "The ID of the PostgreSQL Flexible Server from which to create this PostgreSQL Flexible Server Firewall Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: firewallrules.ValidateFlexibleServerID,
		},

		"end_ip_address": {
			Description:  "The IPv4 Address defining the end of the range of addresses associated with this PostgreSQL Flexible Server Firewall Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsIPv4Address,
		},

		"start_ip_address": {
			Description:  "The IPv4 Address defining the start of the range of addresses associated with this PostgreSQL Flexible Server Firewall Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsIPv4Address,
		},
	}
}
