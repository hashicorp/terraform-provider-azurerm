// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package networkfunction_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/collectorpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkFunctionCollectorPolicyResource struct{}

func TestAccNetworkFunctionCollectorPolicy(t *testing.T) {
	if os.Getenv("ARM_NETWORK_FUNCTION_PEERING_LOCATION") == "" {
		t.Skip("Skipping as ARM_NETWORK_FUNCTION_PEERING_LOCATION is not specified")
	}

	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"Resource": {
			"basic":          testAccNetworkFunctionCollectorPolicy_basic,
			"requiresImport": testAccNetworkFunctionCollectorPolicy_requiresImport,
			"complete":       testAccNetworkFunctionCollectorPolicy_complete,
			"update":         testAccNetworkFunctionCollectorPolicy_update,
		},
	})
}

func testAccNetworkFunctionCollectorPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_collector_policy", "test")
	r := NetworkFunctionCollectorPolicyResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkFunctionCollectorPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_collector_policy", "test")
	r := NetworkFunctionCollectorPolicyResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccNetworkFunctionCollectorPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_collector_policy", "test")
	r := NetworkFunctionCollectorPolicyResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkFunctionCollectorPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_collector_policy", "test")
	r := NetworkFunctionCollectorPolicyResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r NetworkFunctionCollectorPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := collectorpolicies.ParseCollectorPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.NetworkFunction.CollectorPoliciesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r NetworkFunctionCollectorPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_express_route_port" "test" {
  name                = "acctest-erp-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "%[3]s"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%[1]d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  express_route_port_id = azurerm_express_route_port.test.id
  bandwidth_in_gbps     = 1

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = azurerm_express_route_circuit.test.name
  resource_group_name           = azurerm_resource_group.test.name
  shared_key                    = "ItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.16.0/30"
  secondary_peer_address_prefix = "192.168.16.0/30"
  vlan_id                       = 100
}

resource "azurerm_network_function_azure_traffic_collector" "test" {
  name                = "acctest-nfatc-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  depends_on = [
    azurerm_express_route_circuit_peering.test
  ]
}
`, data.RandomInteger, data.Locations.Primary, os.Getenv("ARM_NETWORK_FUNCTION_PEERING_LOCATION"))
}

func (r NetworkFunctionCollectorPolicyResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_function_collector_policy" "test" {
  name                 = "acctest-nfcp-%d"
  location             = "%s"
  traffic_collector_id = azurerm_network_function_azure_traffic_collector.test.id

  ipfx_emission {
    destination_types = ["AzureMonitor"]
  }

  ipfx_ingestion {
    source_resource_ids = [azurerm_express_route_circuit.test.id]
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r NetworkFunctionCollectorPolicyResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_function_collector_policy" "import" {
  name                 = azurerm_network_function_collector_policy.test.name
  location             = "%s"
  traffic_collector_id = azurerm_network_function_azure_traffic_collector.test.id

  ipfx_emission {
    destination_types = ["AzureMonitor"]
  }

  ipfx_ingestion {
    source_resource_ids = [azurerm_express_route_circuit.test.id]
  }
}
`, config, data.Locations.Primary)
}

func (r NetworkFunctionCollectorPolicyResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_function_collector_policy" "test" {
  name                 = "acctest-nfcp-%d"
  location             = "%s"
  traffic_collector_id = azurerm_network_function_azure_traffic_collector.test.id

  ipfx_emission {
    destination_types = ["AzureMonitor"]
  }

  ipfx_ingestion {
    source_resource_ids = [azurerm_express_route_circuit.test.id]
  }

  tags = {
    key = "value"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r NetworkFunctionCollectorPolicyResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_function_collector_policy" "test" {
  name                 = "acctest-nfcp-%d"
  location             = "%s"
  traffic_collector_id = azurerm_network_function_azure_traffic_collector.test.id

  ipfx_emission {
    destination_types = ["AzureMonitor"]
  }

  ipfx_ingestion {
    source_resource_ids = [azurerm_express_route_circuit.test.id]
  }

  tags = {
    key = "value2"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}
