package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ExpressRouteConnectionResource struct{}

func TestAccExpressRouteConnection(t *testing.T) {
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"Resource": {
			"basic":          testAccExpressRouteConnection_basic,
			"requiresImport": testAccExpressRouteConnection_requiresImport,
			"complete":       testAccExpressRouteConnection_complete,
			"update":         testAccExpressRouteConnection_update,
		},
	})
}

func testAccExpressRouteConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_connection", "test")
	r := ExpressRouteConnectionResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_express_route_connection.test").Key("routing.0.associated_route_table_id").Exists(),
				check.That("azurerm_express_route_connection.test").Key("routing.0.propagated_route_table.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func testAccExpressRouteConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_connection", "test")
	r := ExpressRouteConnectionResource{}

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

func testAccExpressRouteConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_connection", "test")
	r := ExpressRouteConnectionResource{}

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

func testAccExpressRouteConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_connection", "test")
	r := ExpressRouteConnectionResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
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

func (r ExpressRouteConnectionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	expressRouteConnectionClient := client.Network.ExpressRouteConnectionsClient
	id, err := parse.ExpressRouteConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := expressRouteConnectionClient.Get(ctx, id.ResourceGroup, id.ExpressRouteGatewayName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving Express Route Connection %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.ExpressRouteConnectionProperties != nil), nil
}

func (r ExpressRouteConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_connection" "test" {
  name                             = "acctest-ExpressRouteConnection-%d"
  express_route_gateway_id         = azurerm_express_route_gateway.test.id
  express_route_circuit_peering_id = azurerm_express_route_circuit_peering.test.id

  depends_on = [azurerm_virtual_hub_route_table.test]
}
`, r.template(data), data.RandomInteger)
}

func (r ExpressRouteConnectionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_connection" "import" {
  name                             = azurerm_express_route_connection.test.name
  express_route_gateway_id         = azurerm_express_route_connection.test.express_route_gateway_id
  express_route_circuit_peering_id = azurerm_express_route_connection.test.express_route_circuit_peering_id
}
`, config)
}

func (r ExpressRouteConnectionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_connection" "test" {
  name                             = "acctest-ExpressRouteConnection-%d"
  express_route_gateway_id         = azurerm_express_route_gateway.test.id
  express_route_circuit_peering_id = azurerm_express_route_circuit_peering.test.id
  routing_weight                   = 2
  authorization_key                = "90f8db47-e25b-4b65-a68b-7743ced2a16b"
  enable_internet_security         = true

  routing {
    associated_route_table_id = azurerm_virtual_hub_route_table.test.id

    propagated_route_table {
      labels          = ["label1"]
      route_table_ids = [azurerm_virtual_hub_route_table.test.id]
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ExpressRouteConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-erconnection-%d"
  location = "%s"
}

resource "azurerm_express_route_port" "test" {
  name                = "acctest-erp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "CDC-Canberra"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  express_route_port_id = azurerm_express_route_port.test.id
  bandwidth_in_gbps     = 5

  sku {
    tier   = "Premium"
    family = "MeteredData"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = azurerm_express_route_circuit.test.name
  resource_group_name           = azurerm_resource_group.test.name
  shared_key                    = "ItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.2.0/30"
  vlan_id                       = 100
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctest-vwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-vhub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"
}

resource "azurerm_express_route_gateway" "test" {
  name                = "acctest-ergw-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_hub_id      = azurerm_virtual_hub.test.id
  scale_units         = 1
}

resource "azurerm_virtual_hub_route_table" "test" {
  name           = "acctest-vhubrt-%d"
  virtual_hub_id = azurerm_virtual_hub.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
