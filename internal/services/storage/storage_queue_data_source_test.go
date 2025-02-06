// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type StorageQueueDataSource struct{}

func TestAccDataSourceStorageQueue_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_queue", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageQueueDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("metadata.%").HasValue("2"),
				check.That(data.ResourceName).Key("metadata.k1").HasValue("v1"),
				check.That(data.ResourceName).Key("metadata.k2").HasValue("v2"),
			),
		},
	})
}

func (d StorageQueueDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "queuedstest-%[1]s"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsadsc%[1]s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_queue" "test" {
  name                 = "queuedstest-%[1]s"
  storage_account_name = "${azurerm_storage_account.test.name}"
  metadata = {
    k1 = "v1"
    k2 = "v2"
  }
}

data "azurerm_storage_queue" "test" {
  name                 = azurerm_storage_queue.test.name
  storage_account_name = azurerm_storage_queue.test.storage_account_name
}
`, data.RandomString, data.Locations.Primary, data.RandomInteger)
}
