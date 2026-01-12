// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
)

type ExascaleDatabasetorageVaultDataSource struct{}

func TestDbStorageVaultDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, fmt.Sprintf("data.%[1]s", oracle.ExascaleDatabaseStorageVaultDataSource{}.ResourceType()), "test")
	r := ExascaleDatabasetorageVaultDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("high_capacity_database_storage.#").HasValue("1"),
			),
		},
	})
}

func (d ExascaleDatabasetorageVaultDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_exascale_database_storage_vault" "test" {
  name                = azurerm_oracle_exascale_database_storage_vault.test.name
  resource_group_name = azurerm_oracle_exascale_database_storage_vault.test.resource_group_name
}
`, ExascaleDatabaseStorageVaultResource{}.basic(data))
}
