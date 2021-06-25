package network_test

import (
	"context"
	"fmt"
	"os"
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
	// NOTE: As provider status of the Express Route Circuit must be set as provisioned manually by service team, so test cases have to use existing Express Route Circuit for testing.
	if os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" || os.Getenv("ARM_TEST_CIRCUIT_NAME") == "" {
		t.Skip("Skipping as ARM_TEST_DATA_RESOURCE_GROUP and/or ARM_TEST_CIRCUIT_NAME are not specified")
		return
	}

	// NOTE: As the provider status of the Express Route Circuit only can be manually updated to `Provisioned` by service team, so currently there is only one authorized test data.
	// And there can be only one Express Route Connection on the same Express Route Circuit, otherwise tests will conflict if run at the same time.
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

resource "azurerm_virtual_hub_route_table" "test" {
  name           = "acctest-VHUBRT-%d"
  virtual_hub_id = azurerm_virtual_hub.test.id
}

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
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ExpressRouteConnectionResource) template(data acceptance.TestData) string {
	circuitRG := os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP")
	circuitName := os.Getenv("ARM_TEST_CIRCUIT_NAME")

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = "%s"
  resource_group_name           = "%s"
  shared_key                    = "ItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.2.0/30"
  vlan_id                       = 100
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-erconnection-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctest-VWAN-%d"
  resource_group_name = azurerm_resource_group.test2.name
  location            = azurerm_resource_group.test2.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-VHUB-%d"
  resource_group_name = azurerm_resource_group.test2.name
  location            = azurerm_resource_group.test2.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"
}

resource "azurerm_express_route_gateway" "test" {
  name                = "acctest-ERGW-%d"
  resource_group_name = azurerm_resource_group.test2.name
  location            = azurerm_resource_group.test2.location
  virtual_hub_id      = azurerm_virtual_hub.test.id
  scale_units         = 1
}
`, circuitName, circuitRG, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
