// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccKustoDatabaseDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kusto_database", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: testAccDataSourceKustoDatabase_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(KustoDatabaseResource{}),
				check.That(data.ResourceName).Key("soft_delete_period").Exists(),
				check.That(data.ResourceName).Key("hot_cache_period").Exists(),
				check.That(data.ResourceName).Key("size").Exists(),
			),
		},
	})
}

func testAccDataSourceKustoDatabase_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kusto_database" "test" {
  name                = azurerm_kusto_database.test.name
  resource_group_name = azurerm_kusto_database.test.resource_group_name
  cluster_name        = azurerm_kusto_database.test.cluster_name
}
`, KustoDatabaseResource{}.complete(data))
}
