package storage_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type StorageAccountDataSource struct{}

func TestAccDataSourceStorageAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.basic(data),
		},
		{
			Config: StorageAccountDataSource{}.basicWithDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("account_tier").HasValue("Standard"),
				check.That(data.ResourceName).Key("account_replication_type").HasValue("LRS"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("production"),
			),
		},
	})
}

func TestAccDataSourceStorageAccount_withWriteLock(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.basicWriteLock(data),
		},
		{
			Config: StorageAccountDataSource{}.basicWriteLockWithDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("account_tier").HasValue("Standard"),
				check.That(data.ResourceName).Key("account_replication_type").HasValue("LRS"),
				check.That(data.ResourceName).Key("primary_connection_string").IsEmpty(),
				check.That(data.ResourceName).Key("secondary_connection_string").IsEmpty(),
				check.That(data.ResourceName).Key("primary_blob_connection_string").IsEmpty(),
				check.That(data.ResourceName).Key("secondary_blob_connection_string").IsEmpty(),
				check.That(data.ResourceName).Key("primary_access_key").IsEmpty(),
				check.That(data.ResourceName).Key("secondary_access_key").IsEmpty(),
			),
		},
	})
}

func (d StorageAccountDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsads%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (d StorageAccountDataSource) basicWriteLock(data acceptance.TestData) string {
	template := d.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_storage_account.test.id
  lock_level = "ReadOnly"
}
`, template, data.RandomInteger)
}

func (d StorageAccountDataSource) basicWithDataSource(data acceptance.TestData) string {
	config := d.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, config)
}

func (d StorageAccountDataSource) basicWriteLockWithDataSource(data acceptance.TestData) string {
	config := d.basicWriteLock(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, config)
}
