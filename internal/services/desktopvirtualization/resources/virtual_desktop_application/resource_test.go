// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_application_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/application"
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
  application_type             = "InBuilt"
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
  source            = "${path.module}/../../testdata/%[2]s"
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

data "azuread_service_principal" "avd" {
  client_id = "9cdead84-a844-4324-93f2-b2e6bb768d07"
}

data "azuread_service_principal" "avd_arm_provider" {
  client_id = "50e95039-b200-4007-bc97-8d5790743a63"
}

resource "azurerm_role_assignment" "avd_storage" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Reader and Data Access"
  principal_id         = data.azuread_service_principal.avd.object_id
}

resource "azurerm_role_assignment" "avd_arm_provider_storage" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Reader and Data Access"
  principal_id         = data.azuread_service_principal.avd_arm_provider.object_id
}

%[4]s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/24"]
}

resource "azurerm_subnet" "test" {
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
    subnet_id                     = azurerm_subnet.test.id
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
  size                = "Standard_F1alds_v7"
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

resource "time_sleep" "test" {
  depends_on = [
    azurerm_virtual_machine_extension.test0,
    azurerm_virtual_machine_extension.test1,
    azurerm_virtual_machine_extension.test2,
    azurerm_role_assignment.avd_storage,
    azurerm_role_assignment.avd_arm_provider_storage
  ]

  create_duration = "2m"
}

resource "azurerm_virtual_desktop_app_attach_package" "test" {
  name                = "acctest-aap-%[1]d"
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
    time_sleep.test
  ]
}

resource "azurerm_virtual_desktop_application_group" "test" {
  name                = "acctest-vdag-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  host_pool_id        = azurerm_virtual_desktop_host_pool.test.id
  type                = "RemoteApp"
}

resource "azurerm_virtual_desktop_application" "test" {
  name                         = "acctest-vdapp-%[1]d"
  application_group_id         = azurerm_virtual_desktop_application_group.test.id
  application_type             = "MsixApplication"
  command_line_argument_policy = "DoNotAllow"
  icon_path                    = "\\\\${azurerm_storage_account.test.name}.file.core.windows.net\\${azurerm_storage_share.test.name}\\${azurerm_storage_share_file.test7.name}"
  msix_package_application_id  = azurerm_virtual_desktop_app_attach_package.test.package_applications[0].app_id
  msix_package_family_name     = azurerm_virtual_desktop_app_attach_package.test.package_family_name
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomString, fileShareConfig, time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC3339))
}
