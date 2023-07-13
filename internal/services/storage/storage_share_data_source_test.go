// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type dataSourceStorageShare struct{}

func TestAccDataSourceStorageShare_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_share", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: dataSourceStorageShare{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("quota").HasValue("120"),
				check.That(data.ResourceName).Key("metadata.%").HasValue("2"),
				check.That(data.ResourceName).Key("metadata.k1").HasValue("v1"),
				check.That(data.ResourceName).Key("metadata.k2").HasValue("v2"),
			),
		},
	})
}

func (d dataSourceStorageShare) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "sharedstest-%s"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsadsc%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_kind             = "FileStorage"
  account_tier             = "Premium"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name                 = "sharedstest-%s"
  storage_account_name = "${azurerm_storage_account.test.name}"
  quota                = 120
  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "rwdl"
      start       = "2019-07-02T09:38:21.0000000Z"
      expiry      = "2019-07-02T10:38:21.0000000Z"
    }
  }
  metadata = {
    k1 = "v1"
    k2 = "v2"
  }
}

data "azurerm_storage_share" "test" {
  name                 = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_share.test.storage_account_name
}
`, data.RandomString, data.Locations.Primary, data.RandomString, data.RandomString)
}
