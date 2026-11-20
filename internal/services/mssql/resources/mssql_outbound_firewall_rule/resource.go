// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_outbound_firewall_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/outboundfirewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlOutboundFirewallRuleCreate,
		Read:   resourceMsSqlOutboundFirewallRuleRead,
		Delete: resourceMsSqlOutboundFirewallRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := outboundfirewallrules.ParseOutboundFirewallRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: mssqlOutboundFirewallRuleSchema(),
	}
}
