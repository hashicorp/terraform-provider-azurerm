package databasemigration_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DatabaseMigrationProjectDataSource struct {
}

func TestAccDatabaseMigrationProjectDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")
	r := DatabaseMigrationProjectDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("source_platform").HasValue("SQL"),
				check.That(data.ResourceName).Key("target_platform").HasValue("SQLDB"),
			),
		},
	})
}

func (DatabaseMigrationProjectDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_database_migration_project" "test" {
  name                = azurerm_database_migration_project.test.name
  service_name        = azurerm_database_migration_project.test.service_name
  resource_group_name = azurerm_database_migration_project.test.resource_group_name
}
`, DatabaseMigrationProjectResource{}.basic(data))
}
