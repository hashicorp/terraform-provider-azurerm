package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskCaching(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskCaching(data, "None"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			), {
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskCaching(data, "ReadOnly"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			), {
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskCaching(data, "ReadWrite"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskCustomSize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				// unset
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_authPassword(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskCustomSize(data, 128),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
			{
				// resize a second time to confirm https://github.com/Azure/azure-rest-api-specs/issues/1906
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskCustomSize(data, 256),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskEphemeral(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskEphemeral(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskDiskEncryptionSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDisk_diskEncryptionSet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskStorageAccountTypeStandardLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskStorageAccountType(data, "Standard_LRS"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskStorageAccountTypeStandardSSDLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskStorageAccountType(data, "StandardSSD_LRS"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskStorageAccountTypePremiumLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskStorageAccountType(data, "Premium_LRS"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func TestAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskWriteAcceleratorEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine_scale_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMWindowsVirtualMachineScaleSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskWriteAcceleratorEnabled(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMWindowsVirtualMachineScaleSetExists(data.ResourceName),
				),
			},
			data.ImportStep(
				"admin_password",
			),
		},
	})
}

func testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskCaching(data acceptance.TestData, caching string) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
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
    caching              = "%s"
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
`, template, caching)
}

func testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskCustomSize(data acceptance.TestData, diskSize int) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
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
`, template, diskSize)
}

func testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDisk_diskEncryptionSetDependencies(data acceptance.TestData) string {
	return fmt.Sprintf(`
locals {
  vm_name = "accVM-%d"
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
`, data.RandomInteger, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDisk_diskEncryptionSetResource(data acceptance.TestData) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDisk_diskEncryptionSetDependencies(data)
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

func testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDisk_diskEncryptionSet(data acceptance.TestData) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDisk_diskEncryptionSetResource(data)
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
    storage_account_type   = "Standard_LRS"
    caching                = "ReadOnly"
    disk_encryption_set_id = azurerm_disk_encryption_set.test.id
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
`, template)
}

func testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskEphemeral(data acceptance.TestData) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_windows_virtual_machine_scale_set" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard_F8s_v2" # has to be this large for ephemeral disks on Windows
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
    caching              = "ReadOnly"

    diff_disk_settings {
      option = "Local"
    }
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
`, template)
}

func testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskStorageAccountType(data acceptance.TestData, storageAccountType string) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
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
    storage_account_type = %q
    caching              = "ReadWrite"
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
`, template, storageAccountType)
}

func testAccAzureRMWindowsVirtualMachineScaleSet_disksOSDiskWriteAcceleratorEnabled(data acceptance.TestData, enabled bool) string {
	template := testAccAzureRMWindowsVirtualMachineScaleSet_template(data)
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
    storage_account_type      = "Premium_LRS"
    caching                   = "None"
    write_accelerator_enabled = %t
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
`, template, enabled)
}
