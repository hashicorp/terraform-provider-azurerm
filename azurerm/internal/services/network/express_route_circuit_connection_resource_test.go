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

type ExpressRouteCircuitConnectionResource struct{}

func TestAccExpressRouteCircuitConnection(t *testing.T) {
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"Resource": {
			"basic":          testAccExpressRouteCircuitConnection_basic,
			"requiresImport": testAccExpressRouteCircuitConnection_requiresImport,
			"complete":       testAccExpressRouteCircuitConnection_complete,
			"update":         testAccExpressRouteCircuitConnection_update,
		},
	})
}

func testAccExpressRouteCircuitConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionResource{}
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

func testAccExpressRouteCircuitConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionResource{}
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

func testAccExpressRouteCircuitConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "846a1918-b7a2-4917-b43c-8c4cdaee006a"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccExpressRouteCircuitConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, "846a1918-b7a2-4917-b43c-8c4cdaee006a"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, "946a1918-b7a2-4917-b43c-8c4cdaee006a"),
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

func (r ExpressRouteCircuitConnectionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ExpressRouteCircuitConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.ExpressRouteCircuitConnectionClient.Get(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName, id.ConnectionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r ExpressRouteCircuitConnectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_connection" "test" {
  name                = "acctest-ExpressRouteCircuitConn-%d"
  peering_id          = azurerm_express_route_circuit_peering.test.id
  peer_peering_id     = azurerm_express_route_circuit_peering.peer_test.id
  address_prefix_ipv4 = "192.169.8.0/29"
}
`, r.template(data), data.RandomInteger)
}

func (r ExpressRouteCircuitConnectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_connection" "import" {
  name                = azurerm_express_route_circuit_connection.test.name
  peering_id          = azurerm_express_route_circuit_connection.test.peering_id
  peer_peering_id     = azurerm_express_route_circuit_connection.test.peer_peering_id
  address_prefix_ipv4 = azurerm_express_route_circuit_connection.test.address_prefix_ipv4
}
`, r.basic(data))
}

func (r ExpressRouteCircuitConnectionResource) complete(data acceptance.TestData, authorizationKey string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_connection" "test" {
  name                = "acctest-ExpressRouteCircuitConn-%d"
  peering_id          = azurerm_express_route_circuit_peering.test.id
  peer_peering_id     = azurerm_express_route_circuit_peering.peer_test.id
  address_prefix_ipv4 = "192.169.8.0/29"
  authorization_key   = "%s"
}
`, r.template(data), data.RandomInteger, authorizationKey)
}

func (r ExpressRouteCircuitConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ercircuitconn-%d"
  location = "%s"
}

resource "azurerm_express_route_port" "test" {
  name                = "acctest-erp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "Equinix-Seattle-SE2"
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
    tier   = "Standard"
    family = "MeteredData"
  }
}

resource "azurerm_express_route_port" "peer_test" {
  name                = "acctest-erp2-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "CDC-Canberra"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
}

resource "azurerm_express_route_circuit" "peer_test" {
  name                  = "acctest-erc2-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  express_route_port_id = azurerm_express_route_port.peer_test.id
  bandwidth_in_gbps     = 5

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
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.1.0/30"
  vlan_id                       = 100
}

resource "azurerm_express_route_circuit_peering" "peer_test" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = azurerm_express_route_circuit.peer_test.name
  resource_group_name           = azurerm_resource_group.test.name
  shared_key                    = "ItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.1.0/30"
  vlan_id                       = 100
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
