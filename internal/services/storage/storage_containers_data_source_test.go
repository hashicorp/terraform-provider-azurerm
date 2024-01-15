// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type storageContainersDataSource struct{}

func TestAccDataSourceStorageContainers_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_containers", "test")
	d := storageContainersDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: d.basic(data, "null"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("containers.#").HasValue("2"),
				check.That(data.ResourceName).Key("containers.0.name").HasValue("test1"),
				check.That(data.ResourceName).Key("containers.0.resource_manager_id").HasValue(
					fmt.Sprintf("/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Storage/storageAccounts/acctestacc%s/blobServices/default/containers/test1",
						data.Client().SubscriptionID, data.RandomInteger, data.RandomString),
				),
				check.That(data.ResourceName).Key("containers.0.data_plane_id").HasValue(
					fmt.Sprintf("https://acctestacc%s.blob.core.windows.net/test1", data.RandomString),
				),
				check.That(data.ResourceName).Key("containers.1.name").HasValue("test2"),
				check.That(data.ResourceName).Key("containers.1.resource_manager_id").HasValue(
					fmt.Sprintf("/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Storage/storageAccounts/acctestacc%s/blobServices/default/containers/test2",
						data.Client().SubscriptionID, data.RandomInteger, data.RandomString),
				),
				check.That(data.ResourceName).Key("containers.1.data_plane_id").HasValue(
					fmt.Sprintf("https://acctestacc%s.blob.core.windows.net/test2", data.RandomString),
				),
			),
		},
	})
}

func TestAccDataSourceStorageContainers_prefix(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_containers", "test")
	d := storageContainersDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: d.basic(data, `"test1"`),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("containers.#").HasValue("1"),
				check.That(data.ResourceName).Key("containers.0.name").HasValue("test1"),
				check.That(data.ResourceName).Key("containers.0.resource_manager_id").HasValue(
					fmt.Sprintf("/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Storage/storageAccounts/acctestacc%s/blobServices/default/containers/test1",
						data.Client().SubscriptionID, data.RandomInteger, data.RandomString),
				),
				check.That(data.ResourceName).Key("containers.0.data_plane_id").HasValue(
					fmt.Sprintf("https://acctestacc%s.blob.core.windows.net/test1", data.RandomString),
				),
			),
		},
	})
}

func (d storageContainersDataSource) basic(data acceptance.TestData, prefix string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                            = "acctestacc%s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "test1" {
  name                  = "test1"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_container" "test2" {
  name                  = "test2"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

data "azurerm_storage_containers" "test" {
  storage_account_id = azurerm_storage_account.test.id
  name_prefix        = %s
  depends_on         = [azurerm_storage_container.test1, azurerm_storage_container.test2]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, prefix)
}
