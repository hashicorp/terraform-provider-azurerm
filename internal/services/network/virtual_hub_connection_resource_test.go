// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualHubConnectionResource struct{}

func TestAccVirtualHubConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

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

func TestAccVirtualHubConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_virtual_hub_connection"),
		},
	})
}

func TestAccVirtualHubConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

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

func TestAccVirtualHubConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubConnection_enableInternetSecurity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.enableInternetSecurity(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.enableInternetSecurity(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubConnection_recreateWithSameConnectionName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	vhubData := data
	vhubData.ResourceName = "azurerm_virtual_hub.test"
	resourceGroupName := fmt.Sprintf("acctestRG-vhub-%d", data.RandomInteger)
	vhubName := fmt.Sprintf("acctest-VHUB-%d", data.RandomInteger)
	vhubConnectionName := fmt.Sprintf("acctestbasicvhubconn-%d", data.RandomInteger)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.template(data),
			Check: acceptance.ComposeTestCheckFunc(
				vhubData.CheckWithClient(checkVirtualHubConnectionDoesNotExist(resourceGroupName, vhubName, vhubConnectionName)),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubConnection_removeRoutingConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withRoutingConfiguration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubConnection_removePropagatedRouteTable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withRoutingConfiguration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.withoutPropagatedRouteTable(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubConnection_removeVnetStaticRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withRoutingConfiguration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.withoutVnetStaticRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubConnection_requiresLocking(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.requiresLocking(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubConnection_updateRoutingConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withRoutingConfiguration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updateRoutingConfiguration(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubConnection_routeMapAndStaticVnetLocalRouteOverrideCriteria(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}
	nameSuffix := randString()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.routeMapAndStaticVnetLocalRouteOverrideCriteria(data, nameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t VirtualHubConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.HubVirtualNetworkConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.HubVirtualNetworkConnectionClient.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Virtual Hub Network Connection (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func checkVirtualHubConnectionDoesNotExist(resourceGroupName, vhubName, vhubConnectionName string) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		if resp, err := clients.Network.HubVirtualNetworkConnectionClient.Get(ctx, resourceGroupName, vhubName, vhubConnectionName); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return fmt.Errorf("Bad: Get on network.HubVirtualNetworkConnectionClient: %+v", err)
		}

		return fmt.Errorf("Bad: Virtual Hub Connection %q (Resource Group %q) still exists", vhubConnectionName, resourceGroupName)
	}
}

func (r VirtualHubConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctestbasicvhubconn-%[2]d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualHubConnectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_hub_connection" "import" {
  name                      = azurerm_virtual_hub_connection.test.name
  virtual_hub_id            = azurerm_virtual_hub_connection.test.virtual_hub_id
  remote_virtual_network_id = azurerm_virtual_hub_connection.test.remote_virtual_network_id
}
`, r.basic(data))
}

func (r VirtualHubConnectionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvirtnet2%[2]d"
  address_space       = ["10.6.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test2" {
  name                = "acctestnsg2%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test2" {
  name                 = "acctestsubnet2%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefixes     = ["10.6.1.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "test2" {
  subnet_id                 = azurerm_subnet.test2.id
  network_security_group_id = azurerm_network_security_group.test2.id
}

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctestvhubconn-%[2]d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
  internet_security_enabled = false
}

resource "azurerm_virtual_hub_connection" "test2" {
  name                      = "acctestvhubconn2-%[2]d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test2.id
  internet_security_enabled = true
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualHubConnectionResource) enableInternetSecurity(data acceptance.TestData, internetSecurityEnabled bool) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctestbasicvhubconn-%[2]d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
  internet_security_enabled = %t
}
`, r.template(data), data.RandomInteger, internetSecurityEnabled)
}

func (VirtualHubConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-VHUB-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.2.0/24"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r VirtualHubConnectionResource) requiresLocking(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  depends_on = [azurerm_virtual_hub.test]
}

resource "azurerm_subnet" "test" {
  # Creating lots of subnets increases the chance of triggering the race condition
  count = 16

  name                 = "acctestsubnet%[1]d-${count.index}"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = [cidrsubnet("10.5.1.0/24", 4, count.index)]

  enforce_private_link_endpoint_network_policies = true
  enforce_private_link_service_network_policies  = true

  service_endpoints = [
    "Microsoft.AzureActiveDirectory",
    "Microsoft.AzureCosmosDB",
    "Microsoft.ContainerRegistry",
    "Microsoft.EventHub",
    "Microsoft.KeyVault",
    "Microsoft.ServiceBus",
    "Microsoft.Sql",
    "Microsoft.Storage",
    "Microsoft.Web",
  ]
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-VHUB-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.2.0/24"
}

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctestbasicvhubconn-%[1]d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r VirtualHubConnectionResource) withRoutingConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctest-vhubconn-%[2]d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id

  routing {
    propagated_route_table {
      labels = ["label1", "label2"]
    }

    static_vnet_route {
      name                = "testvnetroute"
      address_prefixes    = ["10.0.3.0/24", "10.0.4.0/24"]
      next_hop_ip_address = "10.0.3.5"
    }

    static_vnet_route {
      name                = "testvnetroute2"
      address_prefixes    = ["10.0.5.0/24"]
      next_hop_ip_address = "10.0.5.5"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualHubConnectionResource) withoutPropagatedRouteTable(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctest-vhubconn-%[2]d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id

  routing {
    static_vnet_route {
      name                = "testvnetroute"
      address_prefixes    = ["10.0.3.0/24"]
      next_hop_ip_address = "10.0.3.5"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualHubConnectionResource) withoutVnetStaticRoute(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctest-vhubconn-%[2]d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id

  routing {
    propagated_route_table {
      labels = ["default"]
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualHubConnectionResource) updateRoutingConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctest-vhubconn-%[2]d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id

  routing {
    propagated_route_table {
      labels = ["label3"]
    }

    static_vnet_route {
      name                = "testvnetroute6"
      address_prefixes    = ["10.0.6.0/24", "10.0.7.0/24"]
      next_hop_ip_address = "10.0.6.5"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualHubConnectionResource) routeMapAndStaticVnetLocalRouteOverrideCriteria(data acceptance.TestData, nameSuffix string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_route_map" "test" {
  name           = "acctestrm-%[2]s"
  virtual_hub_id = azurerm_virtual_hub.test.id

  rule {
    name                 = "rule1"
    next_step_if_matched = "Continue"

    action {
      type = "Add"

      parameter {
        as_path = ["22334"]
      }
    }

    match_criterion {
      match_condition = "Contains"
      route_prefix    = ["10.0.0.0/8"]
    }
  }
}

resource "azurerm_route_map" "test2" {
  name           = "acctestrmn-%[2]s"
  virtual_hub_id = azurerm_virtual_hub.test.id

  rule {
    name                 = "rule1"
    next_step_if_matched = "Continue"

    action {
      type = "Add"

      parameter {
        as_path = ["22334"]
      }
    }

    match_criterion {
      match_condition = "Contains"
      route_prefix    = ["10.0.0.0/8"]
    }
  }
}

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctest-vhubconn-%[3]d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id

  routing {
    inbound_route_map_id                      = azurerm_route_map.test.id
    outbound_route_map_id                     = azurerm_route_map.test2.id
    static_vnet_local_route_override_criteria = "Equal"

    propagated_route_table {
      labels = ["label3"]
    }

    static_vnet_route {
      name                = "testvnetroute6"
      address_prefixes    = ["10.0.6.0/24", "10.0.7.0/24"]
      next_hop_ip_address = "10.0.6.5"
    }
  }
}
`, r.template(data), nameSuffix, data.RandomInteger)
}
