package desktopvirtualization_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type VirtualDesktopWorkspaceApplicationGroupAssociationResource struct {
}

func TestAccVirtualDesktopWorkspaceApplicationGroupAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_workspace_application_group_association", "test")
	r := VirtualDesktopWorkspaceApplicationGroupAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualDesktopWorkspaceApplicationGroupAssociation_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_workspace_application_group_association", "test")
	r := VirtualDesktopWorkspaceApplicationGroupAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualDesktopWorkspaceApplicationGroupAssociation_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_workspace_application_group_association", "test")
	r := VirtualDesktopWorkspaceApplicationGroupAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualDesktopWorkspaceApplicationGroupAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_workspace_application_group_association", "test")
	r := VirtualDesktopWorkspaceApplicationGroupAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t VirtualDesktopWorkspaceApplicationGroupAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.WorkspaceApplicationGroupAssociationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DesktopVirtualization.WorkspacesClient.Get(ctx, id.Workspace.ResourceGroup, id.Workspace.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Virtual Desktop Workspace Application Group Association %q (Resource Group: %q) does not exist", id.Workspace.Name, id.Workspace.ResourceGroup)
	}

	if resp.StatusCode == http.StatusNotFound {
		return utils.Bool(false), nil
	}

	if resp.WorkspaceProperties == nil || resp.ApplicationGroupReferences == nil {
		return utils.Bool(false), nil
	}

	for _, app := range *resp.ApplicationGroupReferences {
		if strings.EqualFold(app, id.ApplicationGroup.ID()) {
			return utils.Bool(true), nil
		}
	}

	return utils.Bool(false), nil
}

func (VirtualDesktopWorkspaceApplicationGroupAssociationResource) basic(data acceptance.TestData) string {
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

func (VirtualDesktopWorkspaceApplicationGroupAssociationResource) complete(data acceptance.TestData) string {
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

func (r VirtualDesktopWorkspaceApplicationGroupAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_desktop_workspace_application_group_association" "import" {
  workspace_id         = azurerm_virtual_desktop_workspace_application_group_association.test.workspace_id
  application_group_id = azurerm_virtual_desktop_workspace_application_group_association.test.application_group_id
}
`, r.basic(data))
}
