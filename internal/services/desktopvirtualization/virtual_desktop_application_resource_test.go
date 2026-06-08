// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/application"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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

func TestAccVirtualDesktopApplication_appAttachPackage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_application", "test")
	r := VirtualDesktopApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appAttachPackage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
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

	return pointer.To(resp.Model != nil), nil
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

func (r VirtualDesktopApplicationResource) appAttachPackage(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_desktop_application_group" "test" {
  name                = "acctest-vdag-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  host_pool_id        = azurerm_virtual_desktop_host_pool.test.id
  type                = "RemoteApp"
}

resource "azurerm_virtual_desktop_application" "test" {
  name                         = "acctest-vdapp-%[2]d"
  application_group_id         = azurerm_virtual_desktop_application_group.test.id
  application_type             = "MsixApplication"
  command_line_argument_policy = "DoNotAllow"
  icon_path                    = "\\\\${azurerm_storage_account.test.name}.file.core.windows.net\\${azurerm_storage_share.test.name}\\${azurerm_storage_share_file.test7.name}"
  msix_package_application_id  = azurerm_virtual_desktop_app_attach_package.test.package_applications[0].app_id
  msix_package_family_name     = azurerm_virtual_desktop_app_attach_package.test.package_family_name
}
`, VirtualDesktopAppAttachPackageResource{}.complete(data), data.RandomInteger)
}
