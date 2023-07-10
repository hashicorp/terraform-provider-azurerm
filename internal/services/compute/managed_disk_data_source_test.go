// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagedDiskDataSource struct{}

func TestAccDataSourceManagedDisk_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_disk", "test")
	r := ManagedDiskDataSource{}

	name := fmt.Sprintf("acctestmanageddisk-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data, name, resourceGroupName),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(name),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(resourceGroupName),
				check.That(data.ResourceName).Key("storage_account_type").HasValue("Premium_LRS"),
				check.That(data.ResourceName).Key("disk_size_gb").HasValue("10"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("acctest"),
				check.That(data.ResourceName).Key("zones.#").HasValue("1"),
				check.That(data.ResourceName).Key("zones.0").HasValue("2"),
			),
		},
	})
}

func TestAccDataSourceManagedDisk_basic_withUltraSSD(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_disk", "test")
	r := ManagedDiskDataSource{}

	name := fmt.Sprintf("acctestmanageddisk-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic_withUltraSSD(data, name, resourceGroupName),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("disk_iops_read_write").HasValue("101"),
				check.That(data.ResourceName).Key("disk_mbps_read_write").HasValue("10"),
			),
		},
	})
}

func TestAccDataSourceManagedDisk_diskAccess(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_disk", "test")
	r := ManagedDiskDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.diskAccess(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("network_access_policy").HasValue("AllowPrivate"),
				check.That(data.ResourceName).Key("disk_access_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceManagedDisk_encryptionSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_disk", "test")
	r := ManagedDiskDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.encryptionSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("encryption_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("encryption_settings.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("encryption_settings.0.disk_encryption_key.#").HasValue("1"),
				check.That(data.ResourceName).Key("encryption_settings.0.disk_encryption_key.0.secret_url").Exists(),
				check.That(data.ResourceName).Key("encryption_settings.0.disk_encryption_key.0.source_vault_id").Exists(),
				check.That(data.ResourceName).Key("encryption_settings.0.key_encryption_key.#").HasValue("1"),
				check.That(data.ResourceName).Key("encryption_settings.0.key_encryption_key.0.key_url").Exists(),
				check.That(data.ResourceName).Key("encryption_settings.0.key_encryption_key.0.source_vault_id").Exists(),
			),
		},
	})
}

func (ManagedDiskDataSource) basic(data acceptance.TestData, name string, resourceGroupName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "%s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"
  zone                 = "2"

  tags = {
    environment = "acctest"
  }
}

data "azurerm_managed_disk" "test" {
  name                = azurerm_managed_disk.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, resourceGroupName, data.Locations.Primary, name)
}

func (ManagedDiskDataSource) basic_withUltraSSD(data acceptance.TestData, name string, resourceGroupName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "%s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "UltraSSD_LRS"
  create_option        = "Empty"
  disk_size_gb         = "4"
  disk_iops_read_write = "101"
  disk_mbps_read_write = "10"
  zone                 = "2"

  tags = {
    environment = "acctest"
  }
}

data "azurerm_managed_disk" "test" {
  name                = azurerm_managed_disk.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, resourceGroupName, data.Locations.Primary, name)
}

func (ManagedDiskDataSource) diskAccess(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_disk_access" "test" {
  name                = "accda%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_managed_disk" "test" {
  name                  = "acctestd-%[2]d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  storage_account_type  = "Standard_LRS"
  create_option         = "Empty"
  disk_size_gb          = "4"
  zone                  = "1"
  network_access_policy = "AllowPrivate"
  disk_access_id        = azurerm_disk_access.test.id
}

data "azurerm_managed_disk" "test" {
  name                = azurerm_managed_disk.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.Locations.Primary, data.RandomInteger)
}

func (ManagedDiskDataSource) encryptionSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults       = false
      purge_soft_delete_on_destroy          = false
      purge_soft_deleted_keys_on_destroy    = false
      purge_soft_deleted_secrets_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"
  sku_name            = "standard"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.object_id}"

    key_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "GetRotationPolicy",
    ]

    secret_permissions = [
      "Delete",
      "Get",
      "Set",
    ]
  }

  enabled_for_disk_encryption = true

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%s"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "EC"
  key_size     = 2048

  key_opts = [
    "sign",
    "verify",
  ]
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestd-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  encryption_settings {
    enabled = true

    disk_encryption_key {
      secret_url      = "${azurerm_key_vault_secret.test.id}"
      source_vault_id = "${azurerm_key_vault.test.id}"
    }

    key_encryption_key {
      key_url         = "${azurerm_key_vault_key.test.id}"
      source_vault_id = "${azurerm_key_vault.test.id}"
    }
  }
}

data "azurerm_managed_disk" "test" {
  name                = azurerm_managed_disk.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString, data.RandomInteger)
}
