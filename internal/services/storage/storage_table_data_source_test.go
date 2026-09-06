// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type StorageTableDataSource struct{}

func TestAccDataSourceStorageTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_table", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageTableDataSource{}.basicWithDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("acl.#").HasValue("1"),
				check.That(data.ResourceName).Key("acl.0.id").HasValue("MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"),
				check.That(data.ResourceName).Key("acl.0.access_policy.0.permissions").HasValue("raud"),
			),
		},
	})
}

func TestAccDataSourceStorageTable_basicAzureADAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_table", "test")

	// This is instruct the test client to use AAD auth.
	t.Setenv("ARM_STORAGE_USE_AZUREAD", "true")

	// This test is sequential due to we set the env above.
	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: StorageTableDataSource{}.basicAzureADAuthWithDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("acl.#").HasValue("1"),
				check.That(data.ResourceName).Key("acl.0.id").HasValue("MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"),
				check.That(data.ResourceName).Key("acl.0.access_policy.0.permissions").HasValue("raud"),
			),
		},
	})
}

func (d StorageTableDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "tabledstest-%s"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctesttedsc%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_table" "test" {
  name               = "tabletesttedsc%s"
  storage_account_id = azurerm_storage_account.test.id

  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "raud"
      start       = "2020-11-26T08:49:37.0000000Z"
      expiry      = "2020-11-27T08:49:37.0000000Z"
    }
  }
}
`, data.RandomString, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (d StorageTableDataSource) basicWithDataSource(data acceptance.TestData) string {
	config := d.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_table" "test" {
  name               = azurerm_storage_table.test.name
  storage_account_id = azurerm_storage_table.test.storage_account_id
}
`, config)
}

func (d StorageTableDataSource) basicAzureADAuthWithDataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  storage_use_azuread = true
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                      = "acctestacc%[3]s"
  resource_group_name       = azurerm_resource_group.test.name
  location                  = azurerm_resource_group.test.location
  account_tier              = "Standard"
  account_replication_type  = "LRS"
  shared_access_key_enabled = false

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_table" "test" {
  name               = "acctestst%[1]d"
  storage_account_id = azurerm_storage_account.test.id
  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "raud"
      start       = "2020-11-26T08:49:37.0000000Z"
      expiry      = "2020-11-27T08:49:37.0000000Z"
    }
  }
}

data "azurerm_storage_table" "test" {
  name               = azurerm_storage_table.test.name
  storage_account_id = azurerm_storage_table.test.storage_account_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
