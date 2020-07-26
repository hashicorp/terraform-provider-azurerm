package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualHub_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualHub_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_hub"),
			},
		},
	})
}

func TestAccAzureRMVirtualHub_routes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_route(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualHub_routeUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualHub_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_tags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMVirtualHubExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VirtualHubClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Virtual Hub not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Virtual Hub %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on network.VirtualHubClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualHubDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VirtualHubClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_hub" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on network.VirtualHubClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMVirtualHub_basic(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHub_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "test" {
  name                = "acctestVHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualHub_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHub_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "import" {
  name                = azurerm_virtual_hub.test.name
  location            = azurerm_virtual_hub.test.location
  resource_group_name = azurerm_virtual_hub.test.resource_group_name
  virtual_wan_id      = azurerm_virtual_hub.test.virtual_wan_id
  address_prefix      = azurerm_virtual_hub.test.address_prefix
}
`, template)
}

func testAccAzureRMVirtualHub_route(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHub_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "test" {
  name                = "acctestVHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"

  route {
    address_prefixes    = ["172.0.1.0/24"]
    next_hop_ip_address = "12.34.56.78"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualHub_routeUpdated(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHub_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "test" {
  name                = "acctestVHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"

  route {
    address_prefixes    = ["172.0.1.0/24"]
    next_hop_ip_address = "87.65.43.21"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualHub_tags(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHub_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "test" {
  name                = "acctestVHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"

  tags = {
    Hello = "World"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualHub_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
