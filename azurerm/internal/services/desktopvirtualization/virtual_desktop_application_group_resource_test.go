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

type VirtualDesktopApplicationGroupResource struct {
}

func TestAccVirtualDesktopApplicationGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application_group", "test")
	r := VirtualDesktopApplicationGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccVirtualDesktopApplicationGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application_group", "test")
	r := VirtualDesktopApplicationGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func TestAccVirtualDesktopApplicationGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application_group", "test")
	r := VirtualDesktopApplicationGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccVirtualDesktopApplicationGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application_group", "test")
	r := VirtualDesktopApplicationGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_virtual_desktop_application_group"),
		},
	})
}

func (VirtualDesktopApplicationGroupResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ApplicationGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DesktopVirtualization.ApplicationGroupsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Virtual Desktop Application Group %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
	}

	return utils.Bool(resp.ApplicationGroupProperties != nil), nil
}

func (VirtualDesktopApplicationGroupResource) basic(data acceptance.TestData) string {
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

func (VirtualDesktopApplicationGroupResource) complete(data acceptance.TestData) string {
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

func (r VirtualDesktopApplicationGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_desktop_application_group" "import" {
  name                = azurerm_virtual_desktop_application_group.test.name
  location            = azurerm_virtual_desktop_application_group.test.location
  resource_group_name = azurerm_virtual_desktop_application_group.test.resource_group_name
  type                = azurerm_virtual_desktop_application_group.test.type
  host_pool_id        = azurerm_virtual_desktop_application_group.test.host_pool_id
}
`, r.basic(data))
}
