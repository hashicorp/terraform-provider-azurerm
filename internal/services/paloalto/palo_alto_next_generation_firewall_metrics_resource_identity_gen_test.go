// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package paloalto_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccPaloAltoNextGenerationFirewallMetrics_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_next_generation_firewall_metrics", "test")
	r := PaloAltoNextGenerationFirewallMetricsResource{}

	checkedFields := map[string]struct{}{
		"firewall_name":       {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_palo_alto_next_generation_firewall_metrics.test", checkedFields),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_palo_alto_next_generation_firewall_metrics.test", tfjsonpath.New("firewall_name"), tfjsonpath.New("firewall_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_palo_alto_next_generation_firewall_metrics.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("firewall_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_palo_alto_next_generation_firewall_metrics.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("firewall_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(true),
		data.ImportBlockWithIDStep(true),
	}, false)
}
