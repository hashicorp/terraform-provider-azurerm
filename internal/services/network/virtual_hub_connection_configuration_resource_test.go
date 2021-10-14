package network_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualHubConnectionConfigurationResource struct {
}

func TestAccVirtualHubRouteConnectionConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection_configuration", "test")
	r := VirtualHubConnectionConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_hub_connection_id"),
	})
}

func TestAccVirtualHubRouteConnectionConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection_configuration", "test")
	r := VirtualHubConnectionConfigurationResource{}
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

func TestAccVirtualHubRouteConnectionConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection_configuration", "test")
	r := VirtualHubConnectionConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_hub_connection_id"),
	})
}

func TestAccVirtualHubRouteConnectionConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_connection_configuration", "test")
	r := VirtualHubConnectionConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_hub_connection_id"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_hub_connection_id"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("virtual_hub_connection_id"),
	})
}

func (t VirtualHubConnectionConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.HubVirtualNetworkConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.HubVirtualNetworkConnectionClient.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Virtual Hub Connection (%s): %+v", id, err)
	}

	// the associated route table is always a non-empty string
	// it defaults to the defaultRouteTable

	routeTableID := parse.NewHubRouteTableID(id.SubscriptionId, id.ResourceGroup, id.VirtualHubName, "defaultRouteTable").ID()

	associatedRouteConfigured := *resp.RoutingConfiguration.AssociatedRouteTable.ID != routeTableID
	propagatedLabelsConfigured := !reflect.DeepEqual(resp.RoutingConfiguration.PropagatedRouteTables.Labels, &[]string{"default"})
	r := network.SubResource{
		ID: utils.String(routeTableID),
	}
	propagatedIDsConfigured := !reflect.DeepEqual(resp.RoutingConfiguration.PropagatedRouteTables.Ids, &[]network.SubResource{r})
	staticRoutesConfigured := len(*resp.RoutingConfiguration.VnetRoutes.StaticRoutes) > 0

	return utils.Bool(resp.ID != nil && (associatedRouteConfigured || propagatedLabelsConfigured || propagatedIDsConfigured || staticRoutesConfigured)), nil
}

func (VirtualHubConnectionConfigurationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-VHUB-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-NSG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-SUBNET-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctest-VWAN-%d"
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

func (r VirtualHubConnectionConfigurationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_virtual_hub_connection" "test" {
  name                      = "example-vhubconn-%d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
}

resource "azurerm_virtual_hub_route_table" "test" {
  name           = "acctest-RouteTable-%d"
  virtual_hub_id = azurerm_virtual_hub.test.id
  labels         = ["Label1"]

  route {
    name              = "example-route"
    destinations_type = "CIDR"
    destinations      = ["10.0.0.0/16"]
    next_hop_type     = "ResourceId"
    next_hop          = azurerm_virtual_hub_connection.test.id
  }
}

resource "azurerm_virtual_hub_connection_configuration" "test" {
  virtual_hub_connection_id = azurerm_virtual_hub_connection.test.id
  associated_route_table_id = azurerm_virtual_hub_route_table.test.id
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r VirtualHubConnectionConfigurationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection_configuration" "import" {
  virtual_hub_connection_id = azurerm_virtual_hub_connection_configuration.test.virtual_hub_connection_id
  associated_route_table_id = azurerm_virtual_hub_connection_configuration.test.associated_route_table_id
}
`, r.basic(data))
}

func (r VirtualHubConnectionConfigurationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "example-vhubconn-%d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
}

resource "azurerm_virtual_hub_route_table" "test" {
  name           = "acctest-RouteTable-%d"
  virtual_hub_id = azurerm_virtual_hub.test.id
  labels         = ["Label1", "Label2", "Label3"]

  route {
    name              = "example-route"
    destinations_type = "CIDR"
    destinations      = ["10.0.0.0/16"]
    next_hop_type     = "ResourceId"
    next_hop          = azurerm_virtual_hub_connection.test.id
  }
}

resource "azurerm_virtual_hub_connection_configuration" "test" {
  virtual_hub_connection_id = azurerm_virtual_hub_connection.test.id
  associated_route_table_id = azurerm_virtual_hub_route_table.test.id

  propagated_route_table {
    labels = ["Label3"]
    route_table_ids = [
      azurerm_virtual_hub_route_table.test.id,
    ]
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
`, r.template(data), data.RandomInteger, data.RandomInteger)
}
