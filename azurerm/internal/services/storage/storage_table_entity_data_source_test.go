package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type StorageTableEntityDataSource struct{}

func TestAccDataSourceStorageTableEntity_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_table_entity", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: StorageTableEntityDataSource{}.basicWithDataSource(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("entity.%").HasValue("1"),
				check.That(data.ResourceName).Key("entity.testkey").HasValue("testval"),
			),
		},
	})
}

func (d StorageTableEntityDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "tableentitydstest-%s"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctesttedsc%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_table" "test" {
  name                 = "tabletesttedsc%s"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_storage_table_entity" "test" {
  storage_account_name = azurerm_storage_account.test.name
  table_name           = azurerm_storage_table.test.name

  partition_key = "testpartition"
  row_key       = "testrow"

  entity = {
    testkey = "testval"
  }
}
`, data.RandomString, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (d StorageTableEntityDataSource) basicWithDataSource(data acceptance.TestData) string {
	config := d.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_table_entity" "test" {
  table_name           = azurerm_storage_table_entity.test.table_name
  storage_account_name = azurerm_storage_table_entity.test.storage_account_name
  partition_key        = azurerm_storage_table_entity.test.partition_key
  row_key              = azurerm_storage_table_entity.test.row_key
}
`, config)
}
