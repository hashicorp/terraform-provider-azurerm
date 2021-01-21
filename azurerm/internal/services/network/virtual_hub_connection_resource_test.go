package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type VirtualHubConnectionResource struct {
}

func TestAccVirtualHubConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubConnection_enableInternetSecurity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.enableInternetSecurity(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.template(data),
			Check: resource.ComposeTestCheckFunc(
				vhubData.CheckWithClient(checkVirtualHubConnectionDoesNotExist(resourceGroupName, vhubName, vhubConnectionName)),
			),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubConnection_removeRoutingConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withRoutingConfiguration(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubConnection_removePropagatedRouteTable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withRoutingConfiguration(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.withoutPropagatedRouteTable(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubConnection_removeVnetStaticRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withRoutingConfiguration(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.withoutVnetStaticRoute(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubConnection_updateRoutingConfiguration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withRoutingConfiguration(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updateRoutingConfiguration(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t VirtualHubConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
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
	return func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
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
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctestbasicvhubconn-%d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualHubConnectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection" "import" {
  name                      = azurerm_virtual_hub_connection.test.name
  virtual_hub_id            = azurerm_virtual_hub_connection.test.virtual_hub_id
  remote_virtual_network_id = azurerm_virtual_hub_connection.test.remote_virtual_network_id
}
`, r.basic(data))
}

func (r VirtualHubConnectionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvirtnet2%d"
  address_space       = ["10.6.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test2" {
  name                = "acctestnsg2%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test2" {
  name                 = "acctestsubnet2%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefixes     = ["10.6.1.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "test2" {
  subnet_id                 = azurerm_subnet.test2.id
  network_security_group_id = azurerm_network_security_group.test2.id
}

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctestvhubconn-%d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
  internet_security_enabled = false
}

resource "azurerm_virtual_hub_connection" "test2" {
  name                      = "acctestvhubconn2-%d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test2.id
  internet_security_enabled = true
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VirtualHubConnectionResource) enableInternetSecurity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctestbasicvhubconn-%d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
  internet_security_enabled = true
}
`, r.template(data), data.RandomInteger)
}

func (VirtualHubConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vhub-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-VHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.2.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VirtualHubConnectionResource) withRoutingConfiguration(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctest-vhubconn-%d"
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
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctest-vhubconn-%d"
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
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctest-vhubconn-%d"
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
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctest-vhubconn-%d"
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
