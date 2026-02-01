// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccNetworkManagerRoutingRule_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_routing_rule", "test")
	r := NetworkManagerRoutingRuleResource{}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_network_manager_routing_rule.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_network_manager_routing_rule.test", tfjsonpath.New("network_manager_name"), tfjsonpath.New("rule_collection_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_network_manager_routing_rule.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("rule_collection_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_network_manager_routing_rule.test", tfjsonpath.New("routing_configuration_name"), tfjsonpath.New("rule_collection_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_network_manager_routing_rule.test", tfjsonpath.New("rule_collection_name"), tfjsonpath.New("rule_collection_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_network_manager_routing_rule.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("rule_collection_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
