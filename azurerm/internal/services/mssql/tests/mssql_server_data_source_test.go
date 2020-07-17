package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMMsSqlServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMsSqlServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMMsSqlServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMsSqlServer_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "version"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "administrator_login"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fully_qualified_domain_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "tags.%"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMsSqlServer_basic(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlServer_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_mssql_server" "test" {
  name                = azurerm_mssql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
}

`, template)
}

func testAccDataSourceAzureRMMsSqlServer_complete(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlServer_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_mssql_server" "test" {
  name                = azurerm_mssql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
}

`, template)
}
