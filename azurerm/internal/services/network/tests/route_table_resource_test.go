package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type RouteTableResource struct {
}

func TestAccRouteTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disable_bgp_route_propagation").HasValue("false"),
				check.That(data.ResourceName).Key("route.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRouteTable_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disable_bgp_route_propagation").HasValue("false"),
				check.That(data.ResourceName).Key("route.#").HasValue("0"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_route_table"),
		},
	})
}

func TestAccRouteTable_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disable_bgp_route_propagation").HasValue("true"),
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRouteTable_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disable_bgp_route_propagation").HasValue("false"),
				check.That(data.ResourceName).Key("route.#").HasValue("0"),
			),
		},
		{
			Config: r.basicAppliance(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disable_bgp_route_propagation").HasValue("false"),
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
			),
		},
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disable_bgp_route_propagation").HasValue("true"),
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
			),
		},
	})
}

func TestAccRouteTable_singleRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.singleRoute(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRouteTable_removeRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// This configuration includes a single explicit route block
			Config: r.singleRoute(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
			),
		},
		{
			// This configuration has no route blocks at all.
			Config: r.noRouteBlocks(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// The route from the first step is preserved because no
				// blocks at all means "ignore existing blocks".
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
			),
		},
		{
			// This configuration sets route to [] explicitly using the
			// attribute syntax.
			Config: r.singleRouteRemoved(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// The route from the first step is now removed, leaving us
				// with no routes at all.
				check.That(data.ResourceName).Key("route.#").HasValue("0"),
			),
		},
	})
}

func TestAccRouteTable_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckRouteTableDisappears(data.ResourceName),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccRouteTable_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				check.That(data.ResourceName).Key("tags.cost_center").HasValue("MSFT"),
			),
		},
		{
			Config: r.withTagsUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("staging"),
			),
		},
	})
}

func TestAccRouteTable_multipleRoutes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.singleRoute(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
				check.That(data.ResourceName).Key("route.0.name").HasValue("route1"),
				check.That(data.ResourceName).Key("route.0.address_prefix").HasValue("10.1.0.0/16"),
				check.That(data.ResourceName).Key("route.0.next_hop_type").HasValue("VnetLocal"),
			),
		},
		{
			Config: r.multipleRoutes(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("route.#").HasValue("2"),
				check.That(data.ResourceName).Key("route.0.name").HasValue("route1"),
				check.That(data.ResourceName).Key("route.0.address_prefix").HasValue("10.1.0.0/16"),
				check.That(data.ResourceName).Key("route.0.next_hop_type").HasValue("VnetLocal"),
				check.That(data.ResourceName).Key("route.1.name").HasValue("route2"),
				check.That(data.ResourceName).Key("route.1.address_prefix").HasValue("10.2.0.0/16"),
				check.That(data.ResourceName).Key("route.1.next_hop_type").HasValue("VnetLocal"),
			),
		},
		data.ImportStep(),
	})
}

func (t RouteTableResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := network.ParseRouteTableID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.RouteTablesClient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Route Table (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func testCheckRouteTableDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for route table: %q", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.RouteTablesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			if !response.WasNotFound(future.Response()) {
				return fmt.Errorf("Error deleting Route Table %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of Route Table %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		return nil
	}
}

func (RouteTableResource) basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r RouteTableResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_table" "import" {
  name                = azurerm_route_table.test.name
  location            = azurerm_route_table.test.location
  resource_group_name = azurerm_route_table.test.resource_group_name
}
`, r.basic(data))
}

func (RouteTableResource) basicAppliance(data acceptance.TestData) string {
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

  route {
    name                   = "route1"
    address_prefix         = "10.1.0.0/16"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "192.168.0.1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteTableResource) complete(data acceptance.TestData) string {
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

  route {
    name           = "acctestRoute"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "vnetlocal"
  }

  disable_bgp_route_propagation = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteTableResource) singleRoute(data acceptance.TestData) string {
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

  route {
    name           = "route1"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "vnetlocal"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteTableResource) noRouteBlocks(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteTableResource) singleRouteRemoved(data acceptance.TestData) string {
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

  route = []
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteTableResource) multipleRoutes(data acceptance.TestData) string {
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

  route {
    name           = "route1"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "vnetlocal"
  }

  route {
    name           = "route2"
    address_prefix = "10.2.0.0/16"
    next_hop_type  = "vnetlocal"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteTableResource) withTags(data acceptance.TestData) string {
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

  route {
    name           = "route1"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "vnetlocal"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteTableResource) withTagsUpdate(data acceptance.TestData) string {
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

  route {
    name           = "route1"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "vnetlocal"
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
