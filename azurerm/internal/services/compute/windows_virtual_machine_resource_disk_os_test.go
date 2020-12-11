package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccWindowsVirtualMachine_diskOSBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_diskOSBasic(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_diskOSCachingType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_diskOSCachingType(data, "None"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_diskOSCachingType(data, "ReadOnly"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_diskOSCachingType(data, "ReadWrite"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_diskOSCustomName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_diskOSCustomName(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_diskOSCustomSize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_diskOSCustomSize(data, 130),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_diskOSCustomSizeExpanded(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_diskOSCustomSize(data, 130),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_diskOSCustomSize(data, 140),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_diskOSDiskEncryptionSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_diskOSDiskDiskEncryptionSetEncrypted(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password"),
		},
	})
}

func TestAccWindowsVirtualMachine_diskOSDiskEncryptionSetUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_diskOSDiskDiskEncryptionSetUnencrypted(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password"),
			{
				Config: testWindowsVirtualMachine_diskOSDiskDiskEncryptionSetEncrypted(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep("admin_password"),
		},
	})
}

func TestAccWindowsVirtualMachine_diskOSEphemeral(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_diskOSEphemeral(data),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_diskOSStorageTypeStandardLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_diskOSStorageAccountType(data, "Standard_LRS"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_diskOSStorageTypeStandardSSDLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_diskOSStorageAccountType(data, "StandardSSD_LRS"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_diskOSStorageTypePremiumLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_diskOSStorageAccountType(data, "Premium_LRS"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_diskOSStorageTypeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testWindowsVirtualMachine_diskOSStorageAccountType(data, "Standard_LRS"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_diskOSStorageAccountType(data, "Premium_LRS"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_diskOSStorageAccountType(data, "StandardSSD_LRS"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testWindowsVirtualMachine_diskOSStorageAccountType(data, "Standard_LRS"),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccWindowsVirtualMachine_diskOSWriteAcceleratorEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: checkWindowsVirtualMachineIsDestroyed,
		Steps: []resource.TestStep{
			{
				// Enabled
				Config: testWindowsVirtualMachine_diskOSWriteAcceleratorEnabled(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// Disabled
				Config: testWindowsVirtualMachine_diskOSWriteAcceleratorEnabled(data, false),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// Enabled
				Config: testWindowsVirtualMachine_diskOSWriteAcceleratorEnabled(data, true),
				Check: resource.ComposeTestCheckFunc(
					checkWindowsVirtualMachineExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func testWindowsVirtualMachine_diskOSBasic(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_diskOSCachingType(data acceptance.TestData, cachingType string) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "%s"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, cachingType)
}

func testWindowsVirtualMachine_diskOSCustomName(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    name                 = "osdisk1"
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_diskOSCustomSize(data acceptance.TestData, size int) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    name                 = "osdisk1"
    caching              = "ReadWrite"
    disk_size_gb         = %d
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, size)
}

func testWindowsVirtualMachine_diskOSDiskDiskEncryptionSetDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

# note: whilst these aren't used in all tests, it saves us redefining these everywhere
locals {
  vm_name = "acctestvm%s"
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                        = "acctestkv%s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "premium"
  soft_delete_enabled         = true
  purge_protection_enabled    = true
  enabled_for_disk_encryption = true
}

resource "azurerm_key_vault_access_policy" "service-principal" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "create",
    "delete",
    "get",
    "update",
  ]

  secret_permissions = [
    "get",
    "delete",
    "set",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "examplekey"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = ["azurerm_key_vault_access_policy.service-principal"]
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomString, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testWindowsVirtualMachine_diskOSDiskDiskEncryptionSetResource(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_diskOSDiskDiskEncryptionSetDependencies(data)
	return fmt.Sprintf(`
%s

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestdes-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "disk-encryption" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "get",
    "wrapkey",
    "unwrapkey",
  ]

  tenant_id = azurerm_disk_encryption_set.test.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.test.identity.0.principal_id
}

resource "azurerm_role_assignment" "disk-encryption-read-keyvault" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_disk_encryption_set.test.identity.0.principal_id
}
`, template, data.RandomInteger)
}

func testWindowsVirtualMachine_diskOSDiskDiskEncryptionSetUnencrypted(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_diskOSDiskDiskEncryptionSetResource(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  depends_on = [
    "azurerm_role_assignment.disk-encryption-read-keyvault",
    "azurerm_key_vault_access_policy.disk-encryption",
  ]
}
`, template)
}

func testWindowsVirtualMachine_diskOSDiskDiskEncryptionSetEncrypted(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_diskOSDiskDiskEncryptionSetResource(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching                = "ReadWrite"
    storage_account_type   = "Standard_LRS"
    disk_encryption_set_id = azurerm_disk_encryption_set.test.id
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }

  depends_on = [
    "azurerm_role_assignment.disk-encryption-read-keyvault",
    "azurerm_key_vault_access_policy.disk-encryption",
  ]
}
`, template)
}

func testWindowsVirtualMachine_diskOSEphemeral(data acceptance.TestData) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  # Message="OS disk of Ephemeral VM with size greater than 32 GB is not allowed for VM size Standard_F2s_v2"
  size           = "Standard_DS3_v2"
  admin_username = "adminuser"
  admin_password = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadOnly"
    storage_account_type = "Standard_LRS"

    diff_disk_settings {
      option = "Local"
    }
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template)
}

func testWindowsVirtualMachine_diskOSStorageAccountType(data acceptance.TestData, accountType string) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2s_v2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "%s"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, accountType)
}

func testWindowsVirtualMachine_diskOSWriteAcceleratorEnabled(data acceptance.TestData, enabled bool) string {
	template := testWindowsVirtualMachine_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_M8ms"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    storage_account_type      = "Premium_LRS"
    caching                   = "None"
    write_accelerator_enabled = %t
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, template, enabled)
}
