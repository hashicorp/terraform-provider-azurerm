// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_static_website_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type StorageAccountResource struct{}

func (r StorageAccountResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
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
