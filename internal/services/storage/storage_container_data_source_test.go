// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type StorageContainerDataSource struct{}

func TestAccDataSourceStorageContainer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_container", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageContainerDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("container_access_type").HasValue("private"),
				check.That(data.ResourceName).Key("has_immutability_policy").HasValue("false"),
				check.That(data.ResourceName).Key("default_encryption_scope").HasValue(fmt.Sprintf("acctestEScontainer%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("encryption_scope_override_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("metadata.%").HasValue("2"),
				check.That(data.ResourceName).Key("metadata.k1").HasValue("v1"),
				check.That(data.ResourceName).Key("metadata.k2").HasValue("v2"),
			),
		},
	})
}

func (d StorageContainerDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "containerdstest-%[1]s"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsadsc%[1]s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestEScontainer%[3]d"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.Storage"
}

resource "azurerm_storage_container" "test" {
  name                              = "containerdstest-%[1]s"
  storage_account_name              = "${azurerm_storage_account.test.name}"
  container_access_type             = "private"
  default_encryption_scope          = azurerm_storage_encryption_scope.test.name
  encryption_scope_override_enabled = true
  metadata = {
    k1 = "v1"
    k2 = "v2"
  }
}

data "azurerm_storage_container" "test" {
  name                 = azurerm_storage_container.test.name
  storage_account_name = azurerm_storage_container.test.storage_account_name
}
`, data.RandomString, data.Locations.Primary, data.RandomInteger)
}
