package databasemigration_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMDatabaseMigrationService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_database_migration_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatabaseMigrationService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Standard_1vCores"),
				),
			},
		},
	})
}

func testAccDataSourceDatabaseMigrationService_basic(data acceptance.TestData) string {
	config := testAccAzureRMDatabaseMigrationService_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_database_migration_service" "test" {
  resource_group_name = azurerm_database_migration_service.test.resource_group_name
  name                = azurerm_database_migration_service.test.name
}
`, config)
}
