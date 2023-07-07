// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package batch_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type BatchApplicationDataSource struct{}

func TestAccBatchApplicationDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_batch_application", "test")
	r := BatchApplicationDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("allow_updates").HasValue("true"),
			),
		},
	})
}

func (BatchApplicationDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                                = "acctestba%s"
  resource_group_name                 = azurerm_resource_group.test.name
  location                            = azurerm_resource_group.test.location
  pool_allocation_mode                = "BatchService"
  storage_account_id                  = azurerm_storage_account.test.id
  storage_account_authentication_mode = "StorageKeys"
}

resource "azurerm_batch_application" "test" {
  name                = "acctestbatchapp-%d"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
}

data "azurerm_batch_application" "test" {
  name                = azurerm_batch_application.test.name
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_batch_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger)
}
