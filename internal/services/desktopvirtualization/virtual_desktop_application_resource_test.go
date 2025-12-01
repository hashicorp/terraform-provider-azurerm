// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/application"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualDesktopApplicationResource struct{}

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

func TestAccVirtualDesktopApplication_msixType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application", "test")
	r := VirtualDesktopApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.msixType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualDesktopApplication_pathShouldBeSpecifiedForInBuiltApp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application", "test")
	r := VirtualDesktopApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.pathShouldBeSpecifiedForInBuiltApp(),
			ExpectError: regexp.MustCompile(`the 'path' property must be specified when 'application_type' is set to 'InBuilt'`),
		},
	})
}

func TestAccVirtualDesktopApplication_pathCannotBeSpecifiedForMsixApp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application", "test")
	r := VirtualDesktopApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.pathCannotBeSpecifiedForMsixApp(),
			ExpectError: regexp.MustCompile(`the 'path' property cannot be specified when 'application_type' is set to 'MsixApplication'`),
		},
	})
}

func (VirtualDesktopApplicationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := application.ParseApplicationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DesktopVirtualization.ApplicationsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
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
  name                = "acctestHP%d"
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
`, data.RandomInteger, data.Locations.Secondary, data.RandomIntOfLength(8), data.RandomIntOfLength(8), data.RandomIntOfLength(8))
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
  name                 = "acctestHP%d"
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
`, data.RandomInteger, data.Locations.Secondary, data.RandomIntOfLength(8), data.RandomIntOfLength(8), data.RandomIntOfLength(8))
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

func (VirtualDesktopApplicationResource) msixType(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktop-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                = "acctestHP%d"
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

resource "azurerm_virtual_desktop_msix_package" "test" {
  name                = "acctestMSIXPackage%d"
  host_pool_name      = azurerm_virtual_desktop_host_pool.test.name
  resource_group_name = azurerm_resource_group.test.name
  image_path          = "\\\\path\\to\\image.vhd"
  last_updated_in_utc = "2021-09-01T00:00:00"

  package_application {
    app_id            = "app-1"
    app_user_model_id = "app-user-model-1"
    description       = "Testing app 1"
    friendly_name     = "I am your friendly neighbourhood testing app 1"
    icon_image_name   = "icon.png"
    raw_icon          = "VGhpcyBpcyBhIHN0cmluZyB0byBoYXNo"
    raw_png           = "VGhpcyBpcyBhIHN0cmluZyB0byBwYWdl"
  }

  package_family_name   = "msix-package-family-1"
  package_name          = "msix-package-1"
  package_relative_path = "path\\to\\package"
  version               = "0.0.0.1"
}

resource "azurerm_virtual_desktop_application" "test" {
  name                         = "acctestAG%d"
  application_group_id         = azurerm_virtual_desktop_application_group.test.id
  command_line_argument_policy = "DoNotAllow"
  application_type             = "MsixApplication"
  msix_package_app_id          = azurerm_virtual_desktop_msix_package.test.package_application[0].app_id
  msix_package_family_name     = azurerm_virtual_desktop_msix_package.test.package_family_name
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomIntOfLength(8), data.RandomIntOfLength(8), data.RandomIntOfLength(8), data.RandomIntOfLength(8))
}

func (VirtualDesktopApplicationResource) pathShouldBeSpecifiedForInBuiltApp() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_virtual_desktop_application" "test" {
  name                         = "acctestAG"
  application_group_id         = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.DesktopVirtualization/applicationGroups/myapplicationgroup"
  command_line_argument_policy = "DoNotAllow"
}
	`
}

func (VirtualDesktopApplicationResource) pathCannotBeSpecifiedForMsixApp() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_virtual_desktop_application" "test" {
  name                         = "acctestAG"
  application_group_id         = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.DesktopVirtualization/applicationGroups/myapplicationgroup"
  command_line_argument_policy = "DoNotAllow"
  path                         = "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
  application_type             = "MsixApplication"
}
	`
}
