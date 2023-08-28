// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databasemigration_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DatabaseMigrationServiceDataSource struct{}

func TestAccDatabaseMigrationServiceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_database_migration_service", "test")
	r := DatabaseMigrationServiceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("sku_name").HasValue("Standard_1vCores"),
			),
		},
	})
}

func (DatabaseMigrationServiceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_database_migration_service" "test" {
  resource_group_name = azurerm_database_migration_service.test.resource_group_name
  name                = azurerm_database_migration_service.test.name
}
`, DatabaseMigrationServiceResource{}.basic(data))
}
