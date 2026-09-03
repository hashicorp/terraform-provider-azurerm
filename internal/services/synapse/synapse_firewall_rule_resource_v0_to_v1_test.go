// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package synapse_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

// TestAccSynapseFirewallRule_V0ToV1_530 tests the state migration path from an `id` with lowercased
// static segments to their canonicalized format. It uses v5.3.0 as the setup version because it is
// the last release where the `id` could have been stored in state with lowercased static segments via an import using a non-canonical ID.
func TestAccSynapseFirewallRule_V0ToV1_530(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_firewall_rule", "test")
	r := SynapseFirewallRuleResource{}

	importedResourceName := data.ResourceName + "-import"

	data.ResourceRegressionAdditionalStepsTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-synapse-%[2]d/providers/Microsoft.Synapse/workspaces/acctestsw%[2]d/firewallRules/FirewallRule%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.basicV0(data),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Ensure import is a no-op to prevent the resource from normalizing the ID due to a combined CreateUpdate func
					plancheck.ExpectResourceAction(importedResourceName, plancheck.ResourceActionNoop),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourcegroups/acctestRG-synapse-%[2]d/providers/microsoft.synapse/workspaces/acctestsw%[2]d/firewallRules/FirewallRule%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.basicImported(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).ExistsInAzure(r),
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-synapse-%[2]d/providers/Microsoft.Synapse/workspaces/acctestsw%[2]d/firewallRules/FirewallRule%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
	}, "5.3.0")
}

func (r SynapseFirewallRuleResource) basicV0(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_firewall_rule" "test-import" {
  name                 = "FirewallRule%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}

removed {
  from = azurerm_synapse_firewall_rule.test
  lifecycle {
    destroy = false
  }
}

import {
  to = azurerm_synapse_firewall_rule.test-import
  id = "/subscriptions/%[3]s/resourcegroups/acctestRG-synapse-%[2]d/providers/microsoft.synapse/workspaces/acctestsw%[2]d/firewallRules/FirewallRule%[2]d"
}
`, r.template(data), data.RandomInteger, data.Subscriptions.Primary)
}

func (r SynapseFirewallRuleResource) basicImported(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_firewall_rule" "test-import" {
  name                 = "FirewallRule%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}
`, r.template(data), data.RandomInteger)
}
