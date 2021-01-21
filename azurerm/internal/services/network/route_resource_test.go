package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type RouteResource struct {
}

func TestAccRoute_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route", "test")
	r := RouteResource{}

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

func TestAccRoute_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route", "test")
	r := RouteResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_route"),
		},
	})
}

func TestAccRoute_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route", "test")
	r := RouteResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("next_hop_type").HasValue("VnetLocal"),
				check.That(data.ResourceName).Key("next_hop_in_ip_address").HasValue(""),
			),
		},
		{
			Config: r.basicAppliance(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("next_hop_type").HasValue("VirtualAppliance"),
				check.That(data.ResourceName).Key("next_hop_in_ip_address").HasValue("192.168.0.1"),
			),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("next_hop_type").HasValue("VnetLocal"),
				check.That(data.ResourceName).Key("next_hop_in_ip_address").HasValue(""),
			),
		},
	})
}

func TestAccRoute_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route", "test")
	r := RouteResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckRouteDisappears("azurerm_route.test"),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccRoute_multipleRoutes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route", "test")
	r := RouteResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},

		{
			Config: r.multipleRoutes(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (t RouteResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resGroup := id.ResourceGroup
	rtName := id.Path["routeTables"]
	routeName := id.Path["routes"]

	resp, err := clients.Network.RoutesClient.Get(ctx, resGroup, rtName, routeName)
	if err != nil {
		return nil, fmt.Errorf("reading Route (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func testCheckRouteDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.RoutesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		rtName := rs.Primary.Attributes["route_table_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for route: %s", name)
		}

		future, err := client.Delete(ctx, resourceGroup, rtName, name)
		if err != nil {
			return fmt.Errorf("Error deleting Route %q (Route Table %q / Resource Group %q): %+v", name, rtName, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of Route %q (Route Table %q / Resource Group %q): %+v", name, rtName, resourceGroup, err)
		}

		return nil
	}
}

func (RouteResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_route" "test" {
  name                = "acctestroute%d"
  resource_group_name = azurerm_resource_group.test.name
  route_table_name    = azurerm_route_table.test.name

  address_prefix = "10.1.0.0/16"
  next_hop_type  = "vnetlocal"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r RouteResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_route" "import" {
  name                = azurerm_route.test.name
  resource_group_name = azurerm_route.test.resource_group_name
  route_table_name    = azurerm_route.test.route_table_name

  address_prefix = azurerm_route.test.address_prefix
  next_hop_type  = azurerm_route.test.next_hop_type
}
`, r.basic(data))
}

func (RouteResource) basicAppliance(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_route" "test" {
  name                = "acctestroute%d"
  resource_group_name = azurerm_resource_group.test.name
  route_table_name    = azurerm_route_table.test.name

  address_prefix         = "10.1.0.0/16"
  next_hop_type          = "VirtualAppliance"
  next_hop_in_ip_address = "192.168.0.1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (RouteResource) multipleRoutes(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_route" "test" {
  name                = "acctestroute%d"
  resource_group_name = azurerm_resource_group.test.name
  route_table_name    = azurerm_route_table.test.name

  address_prefix = "10.1.0.0/16"
  next_hop_type  = "vnetlocal"
}

resource "azurerm_route" "test1" {
  name                = "acctestroute%d1"
  resource_group_name = azurerm_resource_group.test.name
  route_table_name    = azurerm_route_table.test.name

  address_prefix = "10.2.0.0/16"
  next_hop_type  = "none"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
