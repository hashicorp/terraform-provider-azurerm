package desktopvirtualization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type VirtualDesktopApplicationResource struct {
}

func TestAccVirtualDesktopApplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application", "test")
	r := VirtualDesktopApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualDesktopApplication_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application", "test")
	r := VirtualDesktopApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualDesktopApplication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application", "test")
	r := VirtualDesktopApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualDesktopApplication_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application", "test")
	r := VirtualDesktopApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_virtual_desktop_application"),
		},
	})
}

func (VirtualDesktopApplicationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ApplicationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DesktopVirtualization.ApplicationsClient.Get(ctx, id.ResourceGroup, id.ApplicationGroupName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Virtual Desktop Application Group %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
	}

	return utils.Bool(resp.ApplicationProperties != nil), nil
}

func (VirtualDesktopApplicationResource) basic(data acceptance.TestData) string {
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
  type                = "RemoteApp"
  host_pool_id        = azurerm_virtual_desktop_host_pool.test.id
}

resource "azurerm_virtual_desktop_application" "test" {
  name                         = "acctestAG%d"
  application_group_id         = azurerm_virtual_desktop_application_group.test.id
  path                         = "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
  command_line_argument_policy = "DoNotAllow"
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomIntOfLength(8), data.RandomIntOfLength(8))
}

func (VirtualDesktopApplicationResource) complete(data acceptance.TestData) string {
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
  type                = "RemoteApp"
  host_pool_id        = azurerm_virtual_desktop_host_pool.test.id
  friendly_name       = "TestAppGroup"
  description         = "Acceptance Test: An application group"
  tags = {
    Purpose = "Acceptance-Testing"
  }
}

resource "azurerm_virtual_desktop_application" "test" {
  name                         = "acctestAG%d"
  application_group_id         = azurerm_virtual_desktop_application_group.test.id
  friendly_name                = "Google Chrome"
  description                  = "Chromium based web browser"
  path                         = "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
  command_line_argument_policy = "DoNotAllow"
  command_line_arguments       = "--incognito"
  show_in_portal               = false
  icon_path                    = "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
  icon_index                   = 1
}

`, data.RandomInteger, data.Locations.Secondary, data.RandomIntOfLength(8), data.RandomIntOfLength(8))
}

func (r VirtualDesktopApplicationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_desktop_application" "import" {
  name                         = azurerm_virtual_desktop_application.test.name
  application_group_id         = azurerm_virtual_desktop_application.test.application_group_id
  path                         = azurerm_virtual_desktop_application.test.path
  command_line_argument_policy = azurerm_virtual_desktop_application.test.command_line_argument_policy

}
`, r.basic(data))
}
