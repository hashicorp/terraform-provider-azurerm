// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskCaching(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskCaching(data, "None"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config: r.disksDataDiskCaching(data, "ReadOnly"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config: r.disksDataDiskCaching(data, "ReadWrite"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskDiskEncryptionSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDisk_diskEncryptionSet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskResizing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// 30GB
			Config: r.disksDataDiskResize(data, 30),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			// 60GB
			Config: r.disksDataDiskResize(data, 60),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskMultiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskRemove(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config: r.authPassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskScaling(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// no disks
			Config: r.authPassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			// one disk
			Config: r.disksDataDiskBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			// two disks
			Config: r.disksDataDiskMultiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			// no disks
			Config: r.authPassword(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskStorageAccountTypeStandardLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountType(data, "Standard_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskStorageAccountTypeStandardSSDLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountType(data, "StandardSSD_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskStorageAccountTypeStandardSSDZRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountTypeWithRestrictedLocation(data, "StandardSSD_ZRS", "westeurope"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskStorageAccountTypePremiumLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountType(data, "Premium_LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskStorageAccountTypePremiumV2LRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountTypePremiumV2LRS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskStorageAccountTypePremiumV2LRSWithIOPSAndMBPS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountTypePremiumV2LRSWithIOPSAndMBPS(data, 3001, 126),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
		{
			Config: r.disksDataDiskStorageAccountTypePremiumV2LRSWithIOPSAndMBPS(data, 3002, 127),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskStorageAccountTypePremiumZRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountTypeWithRestrictedLocation(data, "Premium_ZRS", "westeurope"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskStorageAccountTypeUltraSSDLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	// Are supported in East US 2, SouthEast Asia, and North Europe, in two availability zones per region
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountTypeUltraSSDLRS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskStorageAccountTypeUltraSSDLRSWithIOPS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	// Are supported in East US 2, SouthEast Asia, and North Europe, in two availability zones per region
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountTypeUltraSSDLRSWithIOPS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskStorageAccountTypeUltraSSDLRSWithMBPS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	// Are supported in East US 2, SouthEast Asia, and North Europe, in two availability zones per region
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountTypeUltraSSDLRSWithMBPS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskStorageAccountTypeUltraSSDLRSWithIOPSAndMBPS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	// Are supported in East US 2, SouthEast Asia, and North Europe, in two availability zones per region
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskStorageAccountTypeUltraSSDLRSWithIOPSAndMBPS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func TestAccWindowsVirtualMachineScaleSet_disksDataDiskWriteAcceleratorEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")
	r := WindowsVirtualMachineScaleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disksDataDiskWriteAcceleratorEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("admin_password"),
	})
}

func (r WindowsVirtualMachineScaleSetResource) disksDataDiskBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
    disk_size_gb         = 10
    lun                  = 10
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data))
}

func (r WindowsVirtualMachineScaleSetResource) disksDataDiskCaching(data acceptance.TestData, caching string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = "Standard_LRS"
    caching              = %q
    disk_size_gb         = 10
    lun                  = 10
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data), caching)
}

func (r WindowsVirtualMachineScaleSetResource) disksDataDiskResize(data acceptance.TestData, diskSizeGb int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
    lun                  = 10
    disk_size_gb         = %d
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data), diskSizeGb)
}

func (WindowsVirtualMachineScaleSetResource) disksDataDisk_diskEncryptionSetDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults    = false
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

locals {
  vm_name = "acctestVM-%d"
}
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                        = "acctestkv%s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "standard"
  purge_protection_enabled    = true
  soft_delete_retention_days  = 7
  enabled_for_disk_encryption = true
}

resource "azurerm_key_vault_access_policy" "service-principal" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Update",
    "GetRotationPolicy",
  ]

  secret_permissions = [
    "Get",
    "Delete",
    "Set",
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
  address_prefixes     = ["10.0.2.0/24"]
}
`, data.RandomInteger, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r WindowsVirtualMachineScaleSetResource) disksDataDisk_diskEncryptionSetResource(data acceptance.TestData) string {
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
    "Get",
    "WrapKey",
    "UnwrapKey",
    "GetRotationPolicy",
  ]

  tenant_id = azurerm_disk_encryption_set.test.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.test.identity.0.principal_id
}

resource "azurerm_role_assignment" "disk-encryption-read-keyvault" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_disk_encryption_set.test.identity.0.principal_id
}
`, r.disksDataDisk_diskEncryptionSetDependencies(data), data.RandomInteger)
}

func (r WindowsVirtualMachineScaleSetResource) disksDataDisk_diskEncryptionSet(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  computer_name_prefix = "acctestVM"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type   = "Standard_LRS"
    caching                = "ReadWrite"
    disk_encryption_set_id = azurerm_disk_encryption_set.test.id
    disk_size_gb           = 10
    lun                    = 10
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  depends_on = [
    "azurerm_role_assignment.disk-encryption-read-keyvault",
    "azurerm_key_vault_access_policy.disk-encryption",
  ]
}
`, r.disksDataDisk_diskEncryptionSetResource(data))
}

func (r WindowsVirtualMachineScaleSetResource) disksDataDiskMultiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
    disk_size_gb         = 10
    lun                  = 10
  }

  data_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
    disk_size_gb         = 10
    lun                  = 20
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data))
}

func (r WindowsVirtualMachineScaleSetResource) disksDataDiskStorageAccountType(data acceptance.TestData, storageAccountType string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2s_v2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = %q
    caching              = "None"
    disk_size_gb         = 10
    lun                  = 10
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data), storageAccountType)
}

func (r WindowsVirtualMachineScaleSetResource) disksDataDiskStorageAccountTypeWithRestrictedLocation(data acceptance.TestData, storageAccountType string, location string) string {
	// Limited regional availability for some storage account type
	data.Locations.Primary = location
	return r.disksDataDiskStorageAccountType(data, storageAccountType)
}

func (r WindowsVirtualMachineScaleSetResource) disksDataDiskStorageAccountTypePremiumV2LRS(data acceptance.TestData) string {
	// Limited regional availability for `PremiumV2_LRS`
	// `PremiumV2_LRS` disks can only be can only be attached to zonal VMs currently
	data.Locations.Primary = "westeurope"
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2s_v2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  zones               = ["1"]

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = "PremiumV2_LRS"
    caching              = "None"
    disk_size_gb         = 10
    lun                  = 10
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data))
}

func (r WindowsVirtualMachineScaleSetResource) disksDataDiskStorageAccountTypePremiumV2LRSWithIOPSAndMBPS(data acceptance.TestData, iops int, mbps int) string {
	// Limited regional availability for `PremiumV2_LRS`
	data.Locations.Primary = "westeurope"
	if !features.FivePointOh() {
		return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2s_v2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  zones               = ["1"]

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type           = "PremiumV2_LRS"
    caching                        = "None"
    disk_size_gb                   = 10
    lun                            = 10
    ultra_ssd_disk_iops_read_write = %d
    ultra_ssd_disk_mbps_read_write = %d
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data), iops, mbps)
	}

	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F2s_v2"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  zones               = ["1"]

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = "PremiumV2_LRS"
    caching              = "None"
    disk_size_gb         = 10
    lun                  = 10
    disk_iops_read_write = %d
    disk_mbps_read_write = %d
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data), iops, mbps)
}

func (r WindowsVirtualMachineScaleSetResource) disksDataDiskStorageAccountTypeUltraSSDLRS(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_D2s_v3"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  zones               = [1, 2, 3]

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = "UltraSSD_LRS"
    caching              = "None"
    disk_size_gb         = 10
    lun                  = 10
  }

  additional_capabilities {
    ultra_ssd_enabled = true
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data))
}

func (r WindowsVirtualMachineScaleSetResource) disksDataDiskStorageAccountTypeUltraSSDLRSWithIOPS(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_D2s_v3"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  zones               = [1, 2, 3]

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = "UltraSSD_LRS"
    caching              = "None"
    disk_size_gb         = 10
    lun                  = 10
    disk_iops_read_write = 101
  }

  additional_capabilities {
    ultra_ssd_enabled = true
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data))
}

func (r WindowsVirtualMachineScaleSetResource) disksDataDiskStorageAccountTypeUltraSSDLRSWithMBPS(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_D2s_v3"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  zones               = [1, 2, 3]

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = "UltraSSD_LRS"
    caching              = "None"
    disk_size_gb         = 10
    lun                  = 10
    disk_mbps_read_write = 11
  }

  additional_capabilities {
    ultra_ssd_enabled = true
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data))
}

func (r WindowsVirtualMachineScaleSetResource) disksDataDiskStorageAccountTypeUltraSSDLRSWithIOPSAndMBPS(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_D2s_v3"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"
  zones               = [1, 2, 3]

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  data_disk {
    storage_account_type = "UltraSSD_LRS"
    caching              = "None"
    disk_size_gb         = 10
    lun                  = 10
    disk_iops_read_write = 101
    disk_mbps_read_write = 11
  }

  additional_capabilities {
    ultra_ssd_enabled = true
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data))
}

func (r WindowsVirtualMachineScaleSetResource) disksDataDiskWriteAcceleratorEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_M8ms"
  instances           = 1
  admin_username      = "adminuser"
  admin_password      = "P@ssword1234!"

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Premium_LRS"
    caching              = "None"
  }

  data_disk {
    storage_account_type      = "Premium_LRS"
    caching                   = "None"
    disk_size_gb              = 10
    lun                       = 10
    write_accelerator_enabled = true
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.template(data))
}
