package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMMsSqlDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMsSqlDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMsSqlDatabase_basic(data acceptance.TestData) string {
	config := testAccAzureRMMsSqlDatabase_basic(data)
	return fmt.Sprintf(`
%s 

data "azurerm_mssql_database" "test" {
  name            = azurerm_mssql_database.test.name
  mssql_server_id = azurerm_sql_server.test.id
}
`, config)
}
