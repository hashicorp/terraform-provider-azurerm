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

func TestAccAzureRMDesktopVirtualizationWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDesktopVirtualizationWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDesktopVirtualizationWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMDesktopVirtualizationWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDesktopVirtualizationWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDesktopVirtualizationWorkspace_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMDesktopVirtualizationWorkspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDesktopVirtualizationWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDesktopVirtualizationWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccAzureRMDesktopVirtualizationWorkspace_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccAzureRMDesktopVirtualizationWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMDesktopVirtualizationWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDesktopVirtualizationWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDesktopVirtualizationWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationWorkspaceExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMDesktopVirtualizationWorkspace_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_desktop_workspace"),
			},
		},
	})
}

func testCheckAzureRMDesktopVirtualizationWorkspaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DesktopVirtualization.WorkspacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.VirtualDesktopWorkspaceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		result, err := client.Get(ctx, id.ResourceGroup, id.Name)

		if err == nil {
			return nil
		}

		if result.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Virtual Desktop Workspace %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return fmt.Errorf("Bad: Get virtualDesktopWorspaceClient: %+v", err)
	}
}

func testCheckAzureRMDesktopVirtualizationWorkspaceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DesktopVirtualization.WorkspacesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_desktop_workspace" {
			continue
		}

		log.Printf("[WARN] azurerm_virtual_desktop_workspace still exists in state file.")

		id, err := parse.VirtualDesktopWorkspaceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		result, err := client.Get(ctx, id.ResourceGroup, id.Name)

		if err == nil {
			return fmt.Errorf("Virtual Desktop Workspace still exists:\n%#v", result)
		}

		if result.StatusCode != http.StatusNotFound {
			return err
		}
	}

	return nil
}

func testAccAzureRMDesktopVirtualizationWorkspace_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktop-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_workspace" "test" {
  name                = "acctWS%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger)
}

func testAccAzureRMDesktopVirtualizationWorkspace_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktop-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_workspace" "test" {
  name                = "acctestWS%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  friendly_name       = "Acceptance Test!"
  description         = "Acceptance Test by creating acctws%d"
}

`, data.RandomInteger, data.Locations.Secondary, data.RandomIntOfLength(8), data.RandomInteger)
}

func testAccAzureRMDesktopVirtualizationWorkspace_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDesktopVirtualizationWorkspace_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_desktop_workspace" "import" {
  name                = azurerm_virtual_desktop_workspace.test.name
  location            = azurerm_virtual_desktop_workspace.test.location
  resource_group_name = azurerm_virtual_desktop_workspace.test.resource_group_name
}
`, template)
}
