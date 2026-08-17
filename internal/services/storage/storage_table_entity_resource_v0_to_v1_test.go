// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

// TestAccStorageTableEntity_V0ToV1_4810 tests the migration from a data plane URL to a resource manager ID for `storage_table_id`
// It uses v4.81.0 as the setup version because it is the last release where the data plane URL format was accepted.
func TestAccStorageTableEntity_V0ToV1_4810(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table_entity", "test")
	r := StorageTableEntityResource{}

	data.ResourceRegressionTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicV0(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("storage_table_id").HasValue(fmt.Sprintf("https://acctestsa%s.table.core.windows.net/Tables('acctestst%d')", data.RandomString, data.RandomInteger)),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_table_id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.Storage/storageAccounts/acctestsa%[3]s/tableServices/default/tables/acctestst%[2]d", data.Subscriptions.Primary, data.RandomInteger, data.RandomString)),
			),
		},
	}, "4.81.0")
}

func (r StorageTableEntityResource) basicV0(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%[1]d"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_storage_table_entity" "test" {
  storage_table_id = azurerm_storage_table.test.id

  partition_key = "test_partition%[1]d"
  row_key       = "test_row%[1]d"
  entity = {
    Foo = "Bar"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
