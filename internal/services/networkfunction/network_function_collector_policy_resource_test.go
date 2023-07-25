package networkfunction_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/collectorpolicies"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkFunctionCollectorPolicyResource struct{}

func TestAccNetworkFunctionCollectorPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_collector_policy", "test")
	r := NetworkFunctionCollectorPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkFunctionCollectorPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_collector_policy", "test")
	r := NetworkFunctionCollectorPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccNetworkFunctionCollectorPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_collector_policy", "test")
	r := NetworkFunctionCollectorPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetworkFunctionCollectorPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_collector_policy", "test")
	r := NetworkFunctionCollectorPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
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

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-law-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_express_route_port" "test" {
  name                = "acctest-erp-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "Airtel-Chennai2-CLS"
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
  peering_type                  = "MicrosoftPeering"
  express_route_circuit_name    = azurerm_express_route_circuit.test.name
  resource_group_name           = azurerm_resource_group.test.name
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.199.0/30"
  secondary_peer_address_prefix = "192.168.200.0/30"
  vlan_id                       = 300

  microsoft_peering_config {
    advertised_public_prefixes = ["123.6.0.0/24"]
  }
}

resource "azurerm_network_function_azure_traffic_collector" "test" {
  name                = "acctest-nfatc-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r NetworkFunctionCollectorPolicyResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_function_collector_policy" "test" {
  name                                        = "acctest-nfcp-%d"
  network_function_azure_traffic_collector_id = azurerm_network_function_azure_traffic_collector.test.id
  location                                    = "%s"
  emission_policy {
    emission_destination {
      destination_type = "AzureMonitor"
    }
  }
  ingestion_policy {
    ingestion_source {
      resource_id = azurerm_express_route_circuit.test.id
      source_type = "Resource"
    }
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r NetworkFunctionCollectorPolicyResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_function_collector_policy" "import" {
  name                                        = azurerm_network_function_collector_policy.test.name
  network_function_azure_traffic_collector_id = azurerm_network_function_azure_traffic_collector.test.id
  location                                    = "%s"
  emission_policy {
    emission_destination {
      destination_type = "AzureMonitor"
    }
  }
  ingestion_policy {
    ingestion_source {
      resource_id = azurerm_express_route_circuit.test.id
      source_type = "Resource"
    }
  }
}
`, config, data.Locations.Primary)
}

func (r NetworkFunctionCollectorPolicyResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_function_collector_policy" "test" {
  name                                        = "acctest-nfcp-%d"
  network_function_azure_traffic_collector_id = azurerm_network_function_azure_traffic_collector.test.id
  location                                    = "%s"
  emission_policy {
    emission_destination {
      destination_type = "AzureMonitor"
    }
  }
  ingestion_policy {
    ingestion_source {
      resource_id = azurerm_express_route_circuit.test.id
      source_type = "Resource"
    }
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
  name                                        = "acctest-nfcp-%d"
  network_function_azure_traffic_collector_id = azurerm_network_function_azure_traffic_collector.test.id
  location                                    = "%s"
  emission_policy {
    emission_destination {
      destination_type = "AzureMonitor"
    }
  }
  ingestion_policy {
    ingestion_source {
      resource_id = azurerm_express_route_circuit.test.id
      source_type = "Resource"
    }
  }
  tags = {
    key = "value2"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}
