package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRoute_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoute_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists("azurerm_route.test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMRoute_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoute_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMRoute_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_route"),
			},
		},
	})
}

func TestAccAzureRMRoute_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoute_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "next_hop_type", "VnetLocal"),
					resource.TestCheckResourceAttr(data.ResourceName, "next_hop_in_ip_address", ""),
				),
			},
			{
				Config: testAccAzureRMRoute_basicAppliance(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "next_hop_type", "VirtualAppliance"),
					resource.TestCheckResourceAttr(data.ResourceName, "next_hop_in_ip_address", "192.168.0.1"),
				),
			},
			{
				Config: testAccAzureRMRoute_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "next_hop_type", "VnetLocal"),
					resource.TestCheckResourceAttr(data.ResourceName, "next_hop_in_ip_address", ""),
				),
			},
		},
	})
}

func TestAccAzureRMRoute_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoute_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists("azurerm_route.test"),
					testCheckAzureRMRouteDisappears("azurerm_route.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMRoute_multipleRoutes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoute_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists("azurerm_route.test"),
				),
			},

			{
				Config: testAccAzureRMRoute_multipleRoutes(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteExists("azurerm_route.test1"),
				),
			},
		},
	})
}

func testCheckAzureRMRouteExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.RoutesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		rtName := rs.Primary.Attributes["route_table_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for route: %q", name)
		}

		resp, err := client.Get(ctx, resourceGroup, rtName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Route %q (resource group: %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on routesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMRouteDisappears(resourceName string) resource.TestCheckFunc {
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

func testCheckAzureRMRouteDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.RoutesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_route" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		rtName := rs.Primary.Attributes["route_table_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, rtName, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Route still exists:\n%#v", resp.RoutePropertiesFormat)
		}
	}

	return nil
}

func testAccAzureRMRoute_basic(data acceptance.TestData) string {
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

func testAccAzureRMRoute_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_route" "import" {
  name                = azurerm_route.test.name
  resource_group_name = azurerm_route.test.resource_group_name
  route_table_name    = azurerm_route.test.route_table_name

  address_prefix = azurerm_route.test.address_prefix
  next_hop_type  = azurerm_route.test.next_hop_type
}
`, testAccAzureRMRoute_basic(data))
}

func testAccAzureRMRoute_basicAppliance(data acceptance.TestData) string {
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

func testAccAzureRMRoute_multipleRoutes(data acceptance.TestData) string {
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
