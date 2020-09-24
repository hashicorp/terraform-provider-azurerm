package tests

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
)

func TestAccAzureRMVirtualDesktopHostPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_host_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDesktopVirtualizationHostPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualDesktopHostPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationHostPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualDesktopHostPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_host_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDesktopVirtualizationHostPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualDesktopHostPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationHostPoolExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualDesktopHostPool_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_desktop_host_pool"),
			},
		},
	})
}

func testCheckAzureRMDesktopVirtualizationHostPoolExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DesktopVirtualization.HostPoolsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.VirtualDesktopHostPoolID(rs.Primary.ID)
		if err != nil {
			return err
		}

		result, err := client.Get(ctx, id.ResourceGroup, id.Name)

		if err == nil {
			return nil
		}

		if result.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Virtual Desktop Host Pool %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return fmt.Errorf("Bad: Get virtualDesktopHostPoolClient: %+v", err)
	}
}

func testCheckAzureRMDesktopVirtualizationHostPoolDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DesktopVirtualization.HostPoolsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_desktop_host_pool" {
			continue
		}

		log.Printf("[WARN] azurerm_virtual_desktop_host_pool still exists in state file.")

		id, err := parse.VirtualDesktopHostPoolID(rs.Primary.ID)
		if err != nil {
			return err
		}

		result, err := client.Get(ctx, id.ResourceGroup, id.Name)

		if err == nil {
			return fmt.Errorf("Virtual Desktop Host Pool still exists:\n%#v", result)
		}

		if result.StatusCode != http.StatusNotFound {
			return err
		}
	}

	return nil
}

func testAccAzureRMVirtualDesktopHostPool_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                             = "acctws%d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  validation_environment           = true
  description                      = "A description"
  type                             = "Shared"
  load_balancer_type               = "BreadthFirst"
  personal_desktop_assignment_type = "Automatic"
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMVirtualDesktopHostPool_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVirtualDesktopHostPool_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_desktop_host_pool" "import" {
  name                             = azurerm_virtual_desktop_host_pool.test.name
  location                         = azurerm_virtual_desktop_host_pool.test.location
  resource_group_name              = azurerm_virtual_desktop_host_pool.test.resource_group_name
  validation_environment           = azurerm_virtual_desktop_host_pool.test.validation_environment
  description                      = azurerm_virtual_desktop_host_pool.test.description
  type                             = azurerm_virtual_desktop_host_pool.test.type
  load_balancer_type               = azurerm_virtual_desktop_host_pool.test.load_balancer_type
  personal_desktop_assignment_type = azurerm_virtual_desktop_host_pool.test.personal_desktop_assignment_type
}
`, template)
}
