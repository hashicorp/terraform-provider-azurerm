package desktopvirtualization_test

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccVirtualDesktopWorkspaceApplicationGroupAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_workspace_application_group_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAccVirtualDesktopWorkspaceApplicationGroupAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDesktopWorkspaceApplicationGroupAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccVirtualDesktopWorkspaceApplicationGroupAssociationExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccVirtualDesktopWorkspaceApplicationGroupAssociation_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_workspace_application_group_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAccVirtualDesktopWorkspaceApplicationGroupAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDesktopWorkspaceApplicationGroupAssociation_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testAccVirtualDesktopWorkspaceApplicationGroupAssociationExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccVirtualDesktopWorkspaceApplicationGroupAssociation_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_workspace_application_group_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAccVirtualDesktopWorkspaceApplicationGroupAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDesktopWorkspaceApplicationGroupAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccVirtualDesktopWorkspaceApplicationGroupAssociationExists(data.ResourceName),
				),
			},
			{
				Config: testAccVirtualDesktopWorkspaceApplicationGroupAssociation_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testAccVirtualDesktopWorkspaceApplicationGroupAssociationExists(data.ResourceName),
				),
			},
			{
				Config: testAccVirtualDesktopWorkspaceApplicationGroupAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccVirtualDesktopWorkspaceApplicationGroupAssociationExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccVirtualDesktopWorkspaceApplicationGroupAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_workspace_application_group_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testAccVirtualDesktopWorkspaceApplicationGroupAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDesktopWorkspaceApplicationGroupAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccVirtualDesktopWorkspaceApplicationGroupAssociationExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccVirtualDesktopWorkspaceApplicationGroupAssociation_requiresImport),
		},
	})
}

func testAccVirtualDesktopWorkspaceApplicationGroupAssociationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DesktopVirtualization.WorkspacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		toBeSplitID := rs.Primary.ID
		splitID := strings.Split(toBeSplitID, "|")
		if len(splitID) != 2 {
			return fmt.Errorf("Expected ID to be in the format {workspaceID}/{applicationGroup} but got %q", toBeSplitID)
		}

		id, err := parse.WorkspaceID(splitID[0])
		if err != nil {
			return err
		}

		result, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err == nil {
			return nil
		}

		if result.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Virtual Desktop Workspace Application Group Association %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		if result.ApplicationGroupReferences == nil {
			return fmt.Errorf("No Virtual Desktop Workspace <==> Application Group Associations exists for Workspace %q (Resource Group: %q)", id.Name, id.ResourceGroup)
		}

		output := make([]string, 0)
		output = append(output, *result.ApplicationGroupReferences...)

		if !utils.SliceContainsValue(output, splitID[1]) {
			return fmt.Errorf("No Virtual Desktop Workspace <==> Application Group Association exists for Workspace %q and Application Group %q (Resource Group: %q)", id.Name, splitID[1], id.ResourceGroup)
		}

		return nil
	}
}

func testAccVirtualDesktopWorkspaceApplicationGroupAssociationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DesktopVirtualization.ApplicationGroupsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_desktop_application_group" {
			continue
		}

		log.Printf("[WARN] azurerm_virtual_desktop_application_group still exists in state file.")

		id, err := parse.ApplicationGroupID(rs.Primary.ID)
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

func testAccVirtualDesktopWorkspaceApplicationGroupAssociation_basic(data acceptance.TestData) string {
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
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                 = "acctestHPPooled%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  validate_environment = true
  type                 = "Pooled"
  load_balancer_type   = "BreadthFirst"
}

resource "azurerm_virtual_desktop_application_group" "test" {
  name                = "acctestAG%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  friendly_name       = "TestAppGroup"
  description         = "Acceptance Test: An application group"
  type                = "Desktop"
  host_pool_id        = azurerm_virtual_desktop_host_pool.test.id
}

resource "azurerm_virtual_desktop_workspace_application_group_association" "test" {
  workspace_id         = azurerm_virtual_desktop_workspace.test.id
  application_group_id = azurerm_virtual_desktop_application_group.test.id
}

`, data.RandomInteger, data.Locations.Secondary, data.RandomIntOfLength(8), data.RandomIntOfLength(8), data.RandomIntOfLength(8))
}

func testAccVirtualDesktopWorkspaceApplicationGroupAssociation_complete(data acceptance.TestData) string {
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
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                 = "acctestHPPooled%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  validate_environment = true
  type                 = "Pooled"
  load_balancer_type   = "BreadthFirst"
}

resource "azurerm_virtual_desktop_application_group" "test" {
  name                = "acctestAG%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  friendly_name       = "TestAppGroup"
  description         = "Acceptance Test: An application group"
  type                = "Desktop"
  host_pool_id        = azurerm_virtual_desktop_host_pool.test.id
}

resource "azurerm_virtual_desktop_workspace_application_group_association" "test" {
  workspace_id         = azurerm_virtual_desktop_workspace.test.id
  application_group_id = azurerm_virtual_desktop_application_group.test.id
}

resource "azurerm_virtual_desktop_host_pool" "personal" {
  name                             = "acctestHP2nd%d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  type                             = "Personal"
  personal_desktop_assignment_type = "Automatic"
  load_balancer_type               = "Persistent"
}

resource "azurerm_virtual_desktop_application_group" "personal" {
  name                = "acctestAG2nd%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  friendly_name       = "TestAppGroup"
  description         = "Acceptance Test: An application group"
  type                = "Desktop"
  host_pool_id        = azurerm_virtual_desktop_host_pool.personal.id
}

resource "azurerm_virtual_desktop_workspace_application_group_association" "personal" {
  workspace_id         = azurerm_virtual_desktop_workspace.test.id
  application_group_id = azurerm_virtual_desktop_application_group.personal.id
}

`, data.RandomInteger, data.Locations.Secondary, data.RandomIntOfLength(8), data.RandomIntOfLength(8), data.RandomIntOfLength(8), data.RandomIntOfLength(8), data.RandomIntOfLength(8))
}

func testAccVirtualDesktopWorkspaceApplicationGroupAssociation_requiresImport(data acceptance.TestData) string {
	template := testAccVirtualDesktopWorkspaceApplicationGroupAssociation_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_desktop_workspace_application_group_association" "import" {
  workspace_id         = azurerm_virtual_desktop_workspace_application_group_association.test.workspace_id
  application_group_id = azurerm_virtual_desktop_workspace_application_group_association.test.application_group_id
}
`, template)
}
