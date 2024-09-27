// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/connectivityconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerConnectivityConfigurationResource struct{}

func testAccNetworkManagerConnectivityConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_connectivity_configuration", "test")
	r := ManagerConnectivityConfigurationResource{}
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

func testAccNetworkManagerConnectivityConfiguration_basicTopologyMesh(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_connectivity_configuration", "test")
	r := ManagerConnectivityConfigurationResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicTopologyMesh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetworkManagerConnectivityConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_connectivity_configuration", "test")
	r := ManagerConnectivityConfigurationResource{}
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

func testAccNetworkManagerConnectivityConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_connectivity_configuration", "test")
	r := ManagerConnectivityConfigurationResource{}
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

func testAccNetworkManagerConnectivityConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_manager_connectivity_configuration", "test")
	r := ManagerConnectivityConfigurationResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ManagerConnectivityConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := connectivityconfigurations.ParseConnectivityConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Network.ConnectivityConfigurations
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ManagerConnectivityConfigurationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-manager-%d"
  location = "%s"
}

data "azurerm_subscription" "current" {
}

resource "azurerm_network_manager" "test" {
  name                = "acctest-nm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity"]
}

resource "azurerm_network_manager_network_group" "test" {
  name               = "acctest-nmng-%d"
  network_manager_id = azurerm_network_manager.test.id
}

resource "azurerm_virtual_network" "test" {
  name                    = "acctest-vnet-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  address_space           = ["10.0.0.0/16"]
  flow_timeout_in_minutes = 10
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ManagerConnectivityConfigurationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_connectivity_configuration" "test" {
  name                  = "acctest-nmcc-%d"
  network_manager_id    = azurerm_network_manager.test.id
  connectivity_topology = "HubAndSpoke"
  applies_to_group {
    group_connectivity = "None"
    network_group_id   = azurerm_network_manager_network_group.test.id
  }
  hub {
    resource_id   = azurerm_virtual_network.test.id
    resource_type = "Microsoft.Network/virtualNetworks"
  }
}
`, template, data.RandomInteger)
}

func (r ManagerConnectivityConfigurationResource) basicTopologyMesh(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_connectivity_configuration" "test" {
  name                  = "acctest-nmcc-%d"
  network_manager_id    = azurerm_network_manager.test.id
  connectivity_topology = "Mesh"
  applies_to_group {
    group_connectivity = "None"
    network_group_id   = azurerm_network_manager_network_group.test.id
  }
}
`, template, data.RandomInteger)
}

func (r ManagerConnectivityConfigurationResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_connectivity_configuration" "import" {
  name                  = azurerm_network_manager_connectivity_configuration.test.name
  network_manager_id    = azurerm_network_manager_connectivity_configuration.test.network_manager_id
  connectivity_topology = azurerm_network_manager_connectivity_configuration.test.connectivity_topology
  applies_to_group {
    group_connectivity = azurerm_network_manager_connectivity_configuration.test.applies_to_group.0.group_connectivity
    network_group_id   = azurerm_network_manager_connectivity_configuration.test.applies_to_group.0.network_group_id
  }
  hub {
    resource_id   = azurerm_network_manager_connectivity_configuration.test.hub.0.resource_id
    resource_type = azurerm_network_manager_connectivity_configuration.test.hub.0.resource_type
  }
}
`, config)
}

func (r ManagerConnectivityConfigurationResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_network_group" "test2" {
  name               = "acctest-nmng2-%d"
  network_manager_id = azurerm_network_manager.test.id
}

resource "azurerm_network_manager_connectivity_configuration" "test" {
  name                            = "acctest-nmcc-%[2]d"
  network_manager_id              = azurerm_network_manager.test.id
  connectivity_topology           = "HubAndSpoke"
  delete_existing_peering_enabled = false
  global_mesh_enabled             = false
  description                     = "test connectivity configuration"
  applies_to_group {
    group_connectivity  = "None"
    network_group_id    = azurerm_network_manager_network_group.test.id
    global_mesh_enabled = false
    use_hub_gateway     = false
  }
  applies_to_group {
    group_connectivity  = "DirectlyConnected"
    network_group_id    = azurerm_network_manager_network_group.test2.id
    global_mesh_enabled = true
    use_hub_gateway     = true
  }
  hub {
    resource_id   = azurerm_virtual_network.test.id
    resource_type = "Microsoft.Network/virtualNetworks"
  }
}
`, template, data.RandomInteger)
}

func (r ManagerConnectivityConfigurationResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_manager_connectivity_configuration" "test" {
  name                  = "acctest-nmcc-%d"
  network_manager_id    = azurerm_network_manager.test.id
  connectivity_topology = "HubAndSpoke"
  description           = "test"
  global_mesh_enabled   = true
  applies_to_group {
    group_connectivity = "DirectlyConnected"
    network_group_id   = azurerm_network_manager_network_group.test.id
  }
  hub {
    resource_id   = azurerm_virtual_network.test.id
    resource_type = "Microsoft.Network/virtualNetworks"
  }
}
`, template, data.RandomInteger)
}
