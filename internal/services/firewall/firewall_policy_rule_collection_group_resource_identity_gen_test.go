// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package firewall_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccFirewallPolicyRuleCollectionGroup_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall_policy_rule_collection_group", "test")
	r := FirewallPolicyRuleCollectionGroupResource{}

	checkedFields := map[string]struct{}{
		"firewall_policy_name": {},
		"name":                 {},
		"resource_group_name":  {},
		"subscription_id":      {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_firewall_policy_rule_collection_group.test", checkedFields),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_firewall_policy_rule_collection_group.test", tfjsonpath.New("firewall_policy_name"), tfjsonpath.New("firewall_policy_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_firewall_policy_rule_collection_group.test", tfjsonpath.New("name"), tfjsonpath.New("firewall_policy_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_firewall_policy_rule_collection_group.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("firewall_policy_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_firewall_policy_rule_collection_group.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("firewall_policy_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
