// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package redis_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestAccRedisFirewallRule_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_firewall_rule", "test")
	r := RedisFirewallRuleResource{}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectIdentityValue("azurerm_redis_firewall_rule.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_redis_firewall_rule.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_redis_firewall_rule.test", tfjsonpath.New("redis_name"), tfjsonpath.New("redis_cache_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_redis_firewall_rule.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
