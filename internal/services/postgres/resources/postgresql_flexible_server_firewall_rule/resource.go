// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_firewall_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/firewallrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name postgresql_flexible_server_firewall_rule -service-package-name postgres -properties "name" -compare-values "subscription_id:server_id,resource_group_name:server_id,flexible_server_name:server_id"

var postgresqlFlexibleServerResourceName = "azurerm_postgresql_flexible_server"

const azurePostgresqlFlexibleServerFirewallRuleResourceName = "azurerm_postgresql_flexible_server_firewall_rule"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePostgresqlFlexibleServerFirewallRuleCreateUpdate,
		Read:   resourcePostgresqlFlexibleServerFirewallRuleRead,
		Update: resourcePostgresqlFlexibleServerFirewallRuleCreateUpdate,
		Delete: resourcePostgresqlFlexibleServerFirewallRuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&firewallrules.FirewallRuleId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&firewallrules.FirewallRuleId{}),
		},

		Schema: postgresqlFlexibleServerFirewallRuleSchema(),
	}
}
