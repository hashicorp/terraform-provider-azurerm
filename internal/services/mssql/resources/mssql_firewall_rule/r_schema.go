// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_firewall_rule

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mssqlFirewallRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"server_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseSqlServerIDInsensitively),
		},

		"start_ip_address": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.IsIPAddress,
				validation.StringIsNotEmpty,
			),
		},

		"end_ip_address": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.IsIPAddress,
				validation.StringIsNotEmpty,
			),
		},
	}
}
