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

type ExpressRouteGatewayResource struct {
}

func TestAccExpressRouteGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_gateway", "test")
	r := ExpressRouteGatewayResource{}

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

func TestAccExpressRouteGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_gateway", "test")
	r := ExpressRouteGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_express_route_gateway"),
		},
	})
}

func TestAccExpressRouteGateway_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_gateway", "test")
	r := ExpressRouteGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scale_units").HasValue("2"),
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

func (t ExpressRouteGatewayResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resGroup := id.ResourceGroup
	name := id.Path["expressRouteGateways"]

	resp, err := clients.Network.ExpressRouteGatewaysClient.Get(ctx, resGroup, name)
	if err != nil {
		return nil, fmt.Errorf("reading Express Route Gateway (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r ExpressRouteGatewayResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_gateway" "test" {
  name                = "acctestER-gateway-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_hub_id      = azurerm_virtual_hub.test.id
  scale_units         = 1
}
`, r.template(data), data.RandomInteger)
}

func (r ExpressRouteGatewayResource) complete(data acceptance.TestData, scaleUnits int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_gateway" "test" {
  name                = "acctestER-gateway-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_hub_id      = azurerm_virtual_hub.test.id
  scale_units         = %d

  tags = {
    Hello = "World"
  }
}
`, r.template(data), data.RandomInteger, scaleUnits)
}

func (r ExpressRouteGatewayResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_express_route_gateway" "import" {
  name                = azurerm_express_route_gateway.test.name
  resource_group_name = azurerm_express_route_gateway.test.resource_group_name
  location            = azurerm_express_route_gateway.test.location
  virtual_hub_id      = azurerm_express_route_gateway.test.virtual_hub_id
  scale_units         = azurerm_express_route_gateway.test.scale_units
}
`, r.basic(data))
}

func (ExpressRouteGatewayResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-express-%d"
  location = "%s"
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
  address_prefix      = "10.0.1.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
