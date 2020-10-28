package desktopvirtualization_test

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

func TestAccAzureRMVirtualDesktopApplicationGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDesktopVirtualizationApplicationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualDesktopApplicationGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationApplicationGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualDesktopApplicationGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDesktopVirtualizationApplicationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualDesktopApplicationGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationApplicationGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualDesktopApplicationGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDesktopVirtualizationApplicationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualDesktopApplicationGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationApplicationGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccAzureRMVirtualDesktopApplicationGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationApplicationGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			{
				Config: testAccAzureRMVirtualDesktopApplicationGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationApplicationGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualDesktopApplicationGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDesktopVirtualizationApplicationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualDesktopApplicationGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationApplicationGroupExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualDesktopApplicationGroup_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_desktop_application_group"),
			},
		},
	})
}

func testCheckAzureRMDesktopVirtualizationApplicationGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DesktopVirtualization.ApplicationGroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.VirtualDesktopApplicationGroupID(rs.Primary.ID)
		if err != nil {
			return err
		}

		result, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err == nil {
			return nil
		}

		if result.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Virtual Desktop Application Group %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return fmt.Errorf("Bad: Get virtualDesktopApplicationGroupClient: %+v", err)
	}
}

func testCheckAzureRMDesktopVirtualizationApplicationGroupDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DesktopVirtualization.ApplicationGroupsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_desktop_application_group" {
			continue
		}

		log.Printf("[WARN] azurerm_virtual_desktop_application_group still exists in state file.")

		id, err := parse.VirtualDesktopApplicationGroupID(rs.Primary.ID)
		if err != nil {
			return err
		}

		result, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err == nil {
			return fmt.Errorf("Virtual Desktop Application Group still exists:\n%#v", result)
		}

		if result.StatusCode != http.StatusNotFound {
			return err
		}
	}

	return nil
}

func testAccAzureRMVirtualDesktopApplicationGroup_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktop-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                = "acctestHP"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  type                = "Pooled"
  load_balancer_type  = "BreadthFirst"
}

resource "azurerm_virtual_desktop_application_group" "test" {
  name                = "acctestAG%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  type                = "Desktop"
  host_pool_id        = azurerm_virtual_desktop_host_pool.test.id
}

`, data.RandomInteger, data.Locations.Secondary, data.RandomIntOfLength(8))
}

func testAccAzureRMVirtualDesktopApplicationGroup_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktop-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                 = "acctestHP"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  validate_environment = true
  description          = "Acceptance Test: A host pool"
  type                 = "Pooled"
  load_balancer_type   = "BreadthFirst"
}

resource "azurerm_virtual_desktop_application_group" "test" {
  name                = "acctestAG%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  type                = "Desktop"
  host_pool_id        = azurerm_virtual_desktop_host_pool.test.id
  friendly_name       = "TestAppGroup"
  description         = "Acceptance Test: An application group"
  tags = {
    Purpose = "Acceptance-Testing"
  }
}

`, data.RandomInteger, data.Locations.Secondary, data.RandomIntOfLength(8))
}

func testAccAzureRMVirtualDesktopApplicationGroup_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVirtualDesktopApplicationGroup_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_desktop_application_group" "import" {
  name                = azurerm_virtual_desktop_application_group.test.name
  location            = azurerm_virtual_desktop_application_group.test.location
  resource_group_name = azurerm_virtual_desktop_application_group.test.resource_group_name
  type                = azurerm_virtual_desktop_application_group.test.type
  host_pool_id        = azurerm_virtual_desktop_application_group.test.host_pool_id
}
`, template)
}
