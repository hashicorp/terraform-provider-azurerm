package databasemigration_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DatabaseMigrationServiceDataSource struct {
}

func TestAccDatabaseMigrationServiceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_database_migration_service", "test")
	r := DatabaseMigrationServiceDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
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
