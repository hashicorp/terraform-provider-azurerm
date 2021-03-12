package network_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ExpressRouteConnectionResource struct{}

func TestAccExpressRouteConnection_basic(t *testing.T) {
	expressRouteCircuitRg, expressRouteCircuitName, err := testGetExpressRouteCircuitResourceGroupAndNameFromSubscription()
	if err != nil {
		t.Skip(fmt.Sprintf("%+v", err))
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_express_route_connection", "test")
	r := ExpressRouteConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, expressRouteCircuitRg, expressRouteCircuitName),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccExpressRouteConnection_requiresImport(t *testing.T) {
	expressRouteCircuitRG, expressRouteCircuitName, err := testGetExpressRouteCircuitResourceGroupAndNameFromSubscription()
	if err != nil {
		t.Skip(fmt.Sprintf("%+v", err))
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_express_route_connection", "test")
	r := ExpressRouteConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, expressRouteCircuitRG, expressRouteCircuitName),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccExpressRouteConnection_complete(t *testing.T) {
	expressRouteCircuitRG, expressRouteCircuitName, err := testGetExpressRouteCircuitResourceGroupAndNameFromSubscription()
	if err != nil {
		t.Skip(fmt.Sprintf("%+v", err))
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_express_route_connection", "test")
	r := ExpressRouteConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data, expressRouteCircuitRG, expressRouteCircuitName),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccExpressRouteConnection_update(t *testing.T) {
	expressRouteCircuitRG, expressRouteCircuitName, err := testGetExpressRouteCircuitResourceGroupAndNameFromSubscription()
	if err != nil {
		t.Skip(fmt.Sprintf("%+v", err))
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_express_route_connection", "test")
	r := ExpressRouteConnectionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, expressRouteCircuitRG, expressRouteCircuitName),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, expressRouteCircuitRG, expressRouteCircuitName),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, expressRouteCircuitRG, expressRouteCircuitName),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ExpressRouteConnectionResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
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

func (r ExpressRouteConnectionResource) basic(data acceptance.TestData, expressRouteCircuitRG string, expressRouteCircuitName string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_connection" "test" {
  name                             = "acctest-ExpressRouteConnection-%d"
  express_route_gateway_id         = azurerm_express_route_gateway.test.id
  express_route_circuit_peering_id = azurerm_express_route_circuit_peering.test.id
}
`, r.template(data, expressRouteCircuitRG, expressRouteCircuitName), data.RandomInteger)
}

func (r ExpressRouteConnectionResource) requiresImport(data acceptance.TestData) string {
	expressRouteCircuitRG, expressRouteCircuitName, err := testGetExpressRouteCircuitResourceGroupAndNameFromSubscription()
	if err != nil {
		return ""
	}

	config := r.basic(data, expressRouteCircuitRG, expressRouteCircuitName)
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_connection" "import" {
  name                             = azurerm_express_route_connection.test.name
  express_route_gateway_id         = azurerm_express_route_connection.test.express_route_gateway_id
  express_route_circuit_peering_id = azurerm_express_route_connection.test.express_route_circuit_peering_id
}
`, config)
}

func (r ExpressRouteConnectionResource) complete(data acceptance.TestData, expressRouteCircuitRG string, expressRouteCircuitName string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_connection" "test" {
  name                             = "acctest-ExpressRouteConnection-%d"
  express_route_gateway_id         = azurerm_express_route_gateway.test.id
  express_route_circuit_peering_id = azurerm_express_route_circuit_peering.test.id
  routing_weight                   = 2
  authorization_key                = "90f8db47-e25b-4b65-a68b-7743ced2a16b"
  enable_internet_security         = true
}
`, r.template(data, expressRouteCircuitRG, expressRouteCircuitName), data.RandomInteger)
}

func (r ExpressRouteConnectionResource) template(data acceptance.TestData, expressRouteCircuitRG string, expressRouteCircuitName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_resource_group" "test" {
  name = "%s"
}

data "azurerm_express_route_circuit" "test" {
  name                = "%s"
  resource_group_name = data.azurerm_resource_group.test.name
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-erconnection-%d"
  location = data.azurerm_resource_group.test.location
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = data.azurerm_express_route_circuit.test.name
  resource_group_name           = data.azurerm_express_route_circuit.test.resource_group_name
  shared_key                    = "ItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.2.0/30"
  vlan_id                       = 100
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
`, expressRouteCircuitRG, expressRouteCircuitName, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testGetExpressRouteCircuitResourceGroupAndNameFromSubscription() (string, string, error) {
	subscription := strings.ToLower(os.Getenv("ARM_SUBSCRIPTION_ID"))
	expressRouteCircuitRG := ""
	expressRouteCircuitName := ""

	// As provider status of the Express Route Circuit must be set as provisioned manually by service team, so test cases have to use existing Express Route Circuit for testing.
	// Once hashicorp gets the provisioned Express Route Circuit, the same logic for their subscriptions will also be added here.
	if strings.HasPrefix(subscription, "67a9759d") || strings.HasPrefix(subscription, "85b3dbca") {
		expressRouteCircuitRG = "xz3-test"
		expressRouteCircuitName = "tf-er"
	} else {
		return expressRouteCircuitRG, expressRouteCircuitName, fmt.Errorf("Skipping since test is not running as one of the two valid subscriptions allowed to run ExpressRouteConnection tests")
	}

	return expressRouteCircuitRG, expressRouteCircuitName, nil
}
