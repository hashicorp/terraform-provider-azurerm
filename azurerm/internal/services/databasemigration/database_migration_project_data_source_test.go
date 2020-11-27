package databasemigration_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDatabaseMigrationProjectDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_database_migration_project", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatabaseMigrationProject_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "source_platform", "SQL"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_platform", "SQLDB"),
				),
			},
		},
	})
}

func testAccDataSourceDatabaseMigrationProject_basic(data acceptance.TestData) string {
	config := testAccAzureRMDatabaseMigrationProject_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_database_migration_project" "test" {
  name                = azurerm_database_migration_project.test.name
  service_name        = azurerm_database_migration_project.test.service_name
  resource_group_name = azurerm_database_migration_project.test.resource_group_name
}
`, config)
}
