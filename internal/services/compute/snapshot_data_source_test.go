// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SnapshotDataSource struct{}

func TestAccDataSourceSnapshot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_snapshot", "snapshot")
	r := SnapshotDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
	})
}

func TestAccDataSourceSnapshot_encryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_snapshot", "snapshot")
	r := SnapshotDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.encryption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
	})
}

func TestAccDataSourceSnapshot_trustedLaunch(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_snapshot", "snapshot")
	r := SnapshotDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.trustedLaunch(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("trusted_launch_enabled").HasValue("true"),
			),
		},
	})
}

func (SnapshotDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestmd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss_%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.test.id
}

data "azurerm_snapshot" "snapshot" {
  name                = azurerm_snapshot.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (SnapshotDataSource) encryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestmd-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

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
      "Purge",
      "Set",
    ]
  }

  enabled_for_disk_encryption = true
}

resource "azurerm_key_vault_key" "test" {
  name         = "generated-certificate"
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
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-sauce"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss_%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.test.id
  disk_size_gb        = "20"

  encryption_settings {
    disk_encryption_key {
      secret_url      = azurerm_key_vault_secret.test.id
      source_vault_id = azurerm_key_vault.test.id
    }

    key_encryption_key {
      key_url         = azurerm_key_vault_key.test.id
      source_vault_id = azurerm_key_vault.test.id
    }
  }
}

data "azurerm_snapshot" "snapshot" {
  name                = azurerm_snapshot.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger)
}

func (SnapshotDataSource) trustedLaunch(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_platform_image" "test" {
  location  = "%[2]s"
  publisher = "Canonical"
  offer     = "0001-com-ubuntu-server-jammy"
  sku       = "22_04-lts-gen2"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_managed_disk" "test" {
  name                   = "acctestd-%[1]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  os_type                = "Linux"
  create_option          = "FromImage"
  image_reference_id     = data.azurerm_platform_image.test.id
  storage_account_type   = "Standard_LRS"
  hyper_v_generation     = "V2"
  trusted_launch_enabled = true
}

resource "azurerm_snapshot" "test" {
  name                = "acctestss_%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.test.id
}

data "azurerm_snapshot" "snapshot" {
  name                = azurerm_snapshot.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}
