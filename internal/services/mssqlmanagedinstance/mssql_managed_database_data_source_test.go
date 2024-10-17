// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MsSqlManagedDatabaseDataSource struct{}

func TestAccDataSourceMsSqlManagedDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_managed_database", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: MsSqlManagedDatabaseDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("managed_instance_id").Exists(),
			),
		},
	})
}

func (d MsSqlManagedDatabaseDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_mssql_managed_database" "test" {
  name                  = azurerm_mssql_managed_database.test.name
  managed_instance_id   = azurerm_mssql_managed_database.test.managed_instance_id
}
`, MsSqlManagedInstanceResource{}.basic(data))
}
