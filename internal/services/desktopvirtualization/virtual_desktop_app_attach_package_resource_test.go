// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/appattachpackage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VirtualDesktopAppAttachPackageResource struct{}

func TestAccVirtualDesktopAppAttachPackage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_app_attach_package", "test")
	r := VirtualDesktopAppAttachPackageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualDesktopAppAttachPackage_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_app_attach_package", "test")
	r := VirtualDesktopAppAttachPackageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualDesktopAppAttachPackage_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_app_attach_package", "test")
	r := VirtualDesktopAppAttachPackageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccVirtualDesktopAppAttachPackage_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_app_attach_package", "test")
	r := VirtualDesktopAppAttachPackageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (VirtualDesktopAppAttachPackageResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := appattachpackage.ParseAppAttachPackageID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DesktopVirtualization.AppAttachPackagesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r VirtualDesktopAppAttachPackageResource) template(data acceptance.TestData) string {
	cimFileNames := []string{
		"objectid_c57e6597-3a9f-4723-8865-3272302f8c12_0",
		"objectid_c57e6597-3a9f-4723-8865-3272302f8c12_1",
		"objectid_c57e6597-3a9f-4723-8865-3272302f8c12_2",
		"region_c57e6597-3a9f-4723-8865-3272302f8c12_0",
		"region_c57e6597-3a9f-4723-8865-3272302f8c12_1",
		"region_c57e6597-3a9f-4723-8865-3272302f8c12_2",
		"xmlNotepad.cim",
		"icon.png",
	}

	fileShareConfig := ""
	for i, cimFileName := range cimFileNames {
		fileShareConfig += fmt.Sprintf(`
resource "azurerm_storage_share_file" "test%[1]d" {
  name              = "%[2]s"
  source            = "${path.module}/testdata/%[2]s"
  storage_share_url = azurerm_storage_share.test.url
}
`, i, cimFileName)
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktop-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestst%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_replication_type = "LRS"
  account_tier             = "Standard"
}

resource "azurerm_storage_share" "test" {
  name               = "acctest-share-%[1]d"
  quota              = 16
  storage_account_id = azurerm_storage_account.test.id
}

%[4]s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/24"]
}

resource "azurerm_subnet" "test0" {
  name                 = "acctest-snet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/28"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctest-nic-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  ip_configuration {
    name                          = "acctest-ipconfig-%[1]d"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test0.id
  }
}

resource "azurerm_nat_gateway" "test" {
  name                = "acctest-ng-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                = "acctest-vdpool-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  load_balancer_type  = "BreadthFirst"
  type                = "Pooled"
}

resource "azurerm_virtual_desktop_host_pool_registration_info" "test" {
  expiration_date = "%[5]s"
  hostpool_id     = azurerm_virtual_desktop_host_pool.test.id
}

resource "azurerm_windows_virtual_machine" "test" {
  name                = "vm-%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  network_interface_ids = [
    azurerm_network_interface.test.id
  ]
  size                = "Standard_F1als_v7"
  admin_password      = "Password1234"
  admin_username      = "adminuser"
  secure_boot_enabled = true
  vtpm_enabled        = true

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  identity {
    type = "SystemAssigned"
  }

  source_image_reference {
    offer     = "office-365"
    publisher = "microsoftwindowsdesktop"
    sku       = "win11-24h2-avd-m365"
    version   = "latest"
  }
}

resource "azurerm_virtual_machine_extension" "test0" {
  name                 = "acctest-vmext-0-%[1]d"
  publisher            = "Microsoft.Azure.Security.WindowsAttestation"
  type                 = "GuestAttestation"
  type_handler_version = "1.0"
  virtual_machine_id   = azurerm_windows_virtual_machine.test.id

  depends_on = [
    azurerm_nat_gateway.test
  ]
}

resource "azurerm_virtual_machine_extension" "test1" {
  name                 = "acctest-vmext-1-%[1]d"
  publisher            = "Microsoft.Powershell"
  type                 = "DSC"
  type_handler_version = "2.83"
  virtual_machine_id   = azurerm_windows_virtual_machine.test.id

  protected_settings = jsonencode({
    properties = {
      registrationInfoToken = azurerm_virtual_desktop_host_pool_registration_info.test.token
    }
  })

  settings = jsonencode({
    modulesUrl            = "https://wvdportalstorageblob.blob.core.windows.net/galleryartifacts/Configuration_01-20-2022.zip"
    configurationFunction = "Configuration.ps1\\AddSessionHost"
    properties = {
      hostPoolName = azurerm_virtual_desktop_host_pool.test.name
      aadJoin      = true
    }
  })

  depends_on = [
    azurerm_nat_gateway.test
  ]
}

resource "azurerm_virtual_machine_extension" "test2" {
  name                 = "acctest-vmext-2-%[1]d"
  publisher            = "Microsoft.Azure.ActiveDirectory"
  type                 = "AADLoginForWindows"
  type_handler_version = "2.2"
  virtual_machine_id   = azurerm_windows_virtual_machine.test.id

  depends_on = [
    azurerm_nat_gateway.test
  ]
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomString, fileShareConfig, time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC3339))
}

func (r VirtualDesktopAppAttachPackageResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_desktop_app_attach_package" "test" {
  name                = "acctest-msix-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  display_name        = "XmlNotepad"
  host_pool_ids = [
    azurerm_virtual_desktop_host_pool.test.id
  ]
  msix_package_name     = "43906ChrisLovett.XmlNotepad_2.9.0.21_neutral_split.scale-100_hndwmj480pefj"
  storage_share_file_id = azurerm_storage_share_file.test6.id

  depends_on = [
    azurerm_virtual_machine_extension.test0,
    azurerm_virtual_machine_extension.test1,
    azurerm_virtual_machine_extension.test2
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualDesktopAppAttachPackageResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_desktop_app_attach_package" "import" {
  name                  = azurerm_virtual_desktop_app_attach_package.test.name
  resource_group_name   = azurerm_virtual_desktop_app_attach_package.test.resource_group_name
  location              = azurerm_virtual_desktop_app_attach_package.test.location
  display_name          = azurerm_virtual_desktop_app_attach_package.test.display_name
  host_pool_ids         = azurerm_virtual_desktop_app_attach_package.test.host_pool_ids
  msix_package_name     = azurerm_virtual_desktop_app_attach_package.test.msix_package_name
  storage_share_file_id = azurerm_virtual_desktop_app_attach_package.test.storage_share_file_id
}
`, r.basic(data))
}

func (r VirtualDesktopAppAttachPackageResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_desktop_app_attach_package" "test" {
  name                = "acctest-msix-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  display_name        = "XmlNotepadComplete"
  host_pool_ids = [
    azurerm_virtual_desktop_host_pool.test.id
  ]
  msix_package_name              = "43906ChrisLovett.XmlNotepad_2.9.0.21_neutral__hndwmj480pefj"
  storage_share_file_id          = azurerm_storage_share_file.test6.id
  health_check_status_on_failure = "DoNotFail"
  register_at_log_on_enabled     = false
  state_enabled                  = true

  tags = {
    Environment = "Production"
    Foo         = "Bar"
  }

  depends_on = [
    azurerm_virtual_machine_extension.test0,
    azurerm_virtual_machine_extension.test1,
    azurerm_virtual_machine_extension.test2
  ]
}
`, r.template(data), data.RandomInteger)
}
