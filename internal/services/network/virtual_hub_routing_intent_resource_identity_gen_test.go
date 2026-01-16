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

func TestAccVirtualHubRoutingIntent_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_routing_intent", "test")
	r := VirtualHubRoutingIntentResource{}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_virtual_hub_routing_intent.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_virtual_hub_routing_intent.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("virtual_hub_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_virtual_hub_routing_intent.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("virtual_hub_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_virtual_hub_routing_intent.test", tfjsonpath.New("virtual_hub_name"), tfjsonpath.New("virtual_hub_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
