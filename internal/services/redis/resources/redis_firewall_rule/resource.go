// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_firewall_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisfirewallrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redis/resources/redis_firewall_rule/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name redis_firewall_rule -service-package-name redis -properties "name,resource_group_name,redis_name:redis_cache_name" -known-values "subscription_id:data.Subscriptions.Primary"

var redisFirewallRuleResourceName = "azurerm_redis_firewall_rule"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourceRedisFirewallRuleCreateUpdate,
		Read:     resourceRedisFirewallRuleRead,
		Update:   resourceRedisFirewallRuleCreateUpdate,
		Delete:   resourceRedisFirewallRuleDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&redisfirewallrules.FirewallRuleId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&redisfirewallrules.FirewallRuleId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.FirewallRuleV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: redisFirewallRuleSchema(),
	}
}
