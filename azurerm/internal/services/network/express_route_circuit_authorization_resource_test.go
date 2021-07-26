package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ExpressRouteCircuitAuthorizationResource struct {
}

func testAccExpressRouteCircuitAuthorization_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_authorization", "test")
	r := ExpressRouteCircuitAuthorizationResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("authorization_key").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func testAccExpressRouteCircuitAuthorization_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_authorization", "test")
	r := ExpressRouteCircuitAuthorizationResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("authorization_key").Exists(),
			),
		},
		{
			Config:      r.requiresImportConfig(data),
			ExpectError: acceptance.RequiresImportError("azurerm_express_route_circuit_authorization"),
		},
	})
}

func testAccExpressRouteCircuitAuthorization_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_authorization", "test1")
	r := ExpressRouteCircuitAuthorizationResource{}
	secondResourceName := "azurerm_express_route_circuit_authorization.test2"

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("authorization_key").Exists(),
				acceptance.TestCheckResourceAttrSet(secondResourceName, "authorization_key"),
			),
		},
	})
}

func (t ExpressRouteCircuitAuthorizationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resGroup := id.ResourceGroup
	circuitName := id.Path["expressRouteCircuits"]
	name := id.Path["authorizations"]

	resp, err := clients.Network.ExpressRouteAuthsClient.Get(ctx, resGroup, circuitName, name)
	if err != nil {
		return nil, fmt.Errorf("reading Express Route Circuit Authorization (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ExpressRouteCircuitAuthorizationResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  allow_classic_operations = false

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_authorization" "test" {
  name                       = "acctestauth%d"
  express_route_circuit_name = azurerm_express_route_circuit.test.name
  resource_group_name        = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ExpressRouteCircuitAuthorizationResource) requiresImportConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_circuit_authorization" "import" {
  name                       = azurerm_express_route_circuit_authorization.test.name
  express_route_circuit_name = azurerm_express_route_circuit_authorization.test.express_route_circuit_name
  resource_group_name        = azurerm_express_route_circuit_authorization.test.resource_group_name
}
`, r.basicConfig(data))
}

func (ExpressRouteCircuitAuthorizationResource) multipleConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest-erc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  service_provider_name = "Equinix"
  peering_location      = "Silicon Valley"
  bandwidth_in_mbps     = 50

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  allow_classic_operations = false

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}

resource "azurerm_express_route_circuit_authorization" "test1" {
  name                       = "acctestauth1%d"
  express_route_circuit_name = azurerm_express_route_circuit.test.name
  resource_group_name        = azurerm_resource_group.test.name
}

resource "azurerm_express_route_circuit_authorization" "test2" {
  name                       = "acctestauth2%d"
  express_route_circuit_name = azurerm_express_route_circuit.test.name
  resource_group_name        = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
