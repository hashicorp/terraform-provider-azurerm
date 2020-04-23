package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRouteTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "disable_bgp_route_propagation", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRouteTable_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "disable_bgp_route_propagation", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "0"),
				),
			},
			{
				Config:      testAccAzureRMRouteTable_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_route_table"),
			},
		},
	})
}

func TestAccAzureRMRouteTable_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteTable_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists("azurerm_route_table.test"),
					resource.TestCheckResourceAttr(data.ResourceName, "disable_bgp_route_propagation", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRouteTable_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists("azurerm_route_table.test"),
					resource.TestCheckResourceAttr(data.ResourceName, "disable_bgp_route_propagation", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "0"),
				),
			},
			{
				Config: testAccAzureRMRouteTable_basicAppliance(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists("azurerm_route_table.test"),
					resource.TestCheckResourceAttr(data.ResourceName, "disable_bgp_route_propagation", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "1"),
				),
			},
			{
				Config: testAccAzureRMRouteTable_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists("azurerm_route_table.test"),
					resource.TestCheckResourceAttr(data.ResourceName, "disable_bgp_route_propagation", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMRouteTable_singleRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteTable_singleRoute(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRouteTable_removeRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				// This configuration includes a single explicit route block
				Config: testAccAzureRMRouteTable_singleRoute(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "1"),
				),
			},
			{
				// This configuration has no route blocks at all.
				Config: testAccAzureRMRouteTable_noRouteBlocks(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(data.ResourceName),
					// The route from the first step is preserved because no
					// blocks at all means "ignore existing blocks".
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "1"),
				),
			},
			{
				// This configuration sets route to [] explicitly using the
				// attribute syntax.
				Config: testAccAzureRMRouteTable_singleRouteRemoved(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(data.ResourceName),
					// The route from the first step is now removed, leaving us
					// with no routes at all.
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMRouteTable_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(data.ResourceName),
					testCheckAzureRMRouteTableDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMRouteTable_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteTable_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: testAccAzureRMRouteTable_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func TestAccAzureRMRouteTable_multipleRoutes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRouteTable_singleRoute(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.0.name", "route1"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.0.address_prefix", "10.1.0.0/16"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.0.next_hop_type", "VnetLocal"),
				),
			},
			{
				Config: testAccAzureRMRouteTable_multipleRoutes(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.0.name", "route1"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.0.address_prefix", "10.1.0.0/16"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.0.next_hop_type", "VnetLocal"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.1.name", "route2"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.1.address_prefix", "10.2.0.0/16"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.1.next_hop_type", "VnetLocal"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMRouteTableExists(resourceName string) resource.TestCheckFunc {
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

		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Route Table %q (resource group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on routeTablesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMRouteTableDisappears(resourceName string) resource.TestCheckFunc {
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

func testCheckAzureRMRouteTableDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.RouteTablesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_route_table" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Route Table still exists:\n%#v", resp.RouteTablePropertiesFormat)
	}

	return nil
}

func testAccAzureRMRouteTable_basic(data acceptance.TestData) string {
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

func testAccAzureRMRouteTable_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_table" "import" {
  name                = azurerm_route_table.test.name
  location            = azurerm_route_table.test.location
  resource_group_name = azurerm_route_table.test.resource_group_name
}
`, testAccAzureRMRouteTable_basic(data))
}

func testAccAzureRMRouteTable_basicAppliance(data acceptance.TestData) string {
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

func testAccAzureRMRouteTable_complete(data acceptance.TestData) string {
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

func testAccAzureRMRouteTable_singleRoute(data acceptance.TestData) string {
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

func testAccAzureRMRouteTable_noRouteBlocks(data acceptance.TestData) string {
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

func testAccAzureRMRouteTable_singleRouteRemoved(data acceptance.TestData) string {
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

func testAccAzureRMRouteTable_multipleRoutes(data acceptance.TestData) string {
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

func testAccAzureRMRouteTable_withTags(data acceptance.TestData) string {
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

func testAccAzureRMRouteTable_withTagsUpdate(data acceptance.TestData) string {
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
