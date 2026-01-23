// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccDatabricksVirtualNetworkPeering_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_virtual_network_peering", "test")
	r := DatabricksVirtualNetworkPeeringResource{}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectIdentityValue("azurerm_databricks_virtual_network_peering.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_databricks_virtual_network_peering.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_databricks_virtual_network_peering.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("workspace_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_databricks_virtual_network_peering.test", tfjsonpath.New("workspace_name"), tfjsonpath.New("workspace_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
