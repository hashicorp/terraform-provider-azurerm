// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_firewall_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/firewallrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name mysql_flexible_server_firewall_rule -service-package-name mysql -properties "name,resource_group_name,flexible_server_name:server_name" -known-values "subscription_id:data.Subscriptions.Primary"

var mysqlFlexibleServerFirewallResourceName = "azurerm_mysql_flexible_server_firewall_rule"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMySqlFlexibleServerFirewallRuleCreateUpdate,
		Read:   resourceMySqlFlexibleServerFirewallRuleRead,
		Update: resourceMySqlFlexibleServerFirewallRuleCreateUpdate,
		Delete: resourceMySqlFlexibleServerFirewallRuleDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&firewallrules.FirewallRuleId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&firewallrules.FirewallRuleId{}),
		},

		Schema: mysqlFlexibleServerFirewallRuleSchema(),
	}
}
