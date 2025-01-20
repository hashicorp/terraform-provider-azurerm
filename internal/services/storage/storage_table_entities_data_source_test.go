// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type StorageTableEntitiesDataSource struct{}

func TestAccDataSourceStorageTableEntities_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_table_entities", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageTableEntitiesDataSource{}.basicWithDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("items.#").HasValue("2"),
				check.That(data.ResourceName).Key("items.0.properties.%").HasValue("2"),
			),
		},
	})
}

func TestAccDataSourceStorageTableEntities_withSelector(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_table_entities", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageTableEntitiesDataSource{}.basicWithDataSourceAndSelector(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("items.#").HasValue("1"),
			),
		},
	})
}

func (d StorageTableEntitiesDataSource) basic(data acceptance.TestData) string {
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

  allow_nested_items_to_be_public = false
}

resource "azurerm_storage_table" "test" {
  name                 = "tabletesttedsc%s"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_storage_table_entity" "test" {
  storage_table_id = azurerm_storage_table.test.id

  partition_key = "testpartition"
  row_key       = "testrow"

  entity = {
    testkey1 = "testval11"
    testkey2 = "testval12"
  }
}

resource "azurerm_storage_table_entity" "test2" {
  storage_table_id = azurerm_storage_table.test.id

  partition_key = "testpartition"
  row_key       = "testrow2"

  entity = {
    testkey1 = "testval21"
    testkey2 = "testval22"
  }
}

resource "azurerm_storage_table_entity" "testselector" {
  storage_table_id = azurerm_storage_table.test.id

  partition_key = "testselectorpartition"
  row_key       = "testrow"

  entity = {
    testkey1     = "testval31"
    testkey2     = "testval32"
    testselector = "testselectorval"
  }
}
`, data.RandomString, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (d StorageTableEntitiesDataSource) basicWithDataSource(data acceptance.TestData) string {
	config := d.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_table_entities" "test" {
  storage_table_id = azurerm_storage_table.test.id
  filter           = "PartitionKey eq 'testpartition'"

  depends_on = [
    azurerm_storage_table_entity.test,
    azurerm_storage_table_entity.test2,
  ]
}
`, config)
}

func (d StorageTableEntitiesDataSource) basicWithDataSourceAndSelector(data acceptance.TestData) string {
	config := d.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_table_entities" "test" {
  storage_table_id = azurerm_storage_table.test.id
  filter           = "PartitionKey eq 'testselectorpartition'"
  select           = ["testselector"]

  depends_on = [
    azurerm_storage_table_entity.test,
    azurerm_storage_table_entity.test2,
    azurerm_storage_table_entity.testselector,
  ]
}
`, config)
}
