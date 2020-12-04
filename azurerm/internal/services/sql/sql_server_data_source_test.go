package sql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceSqlServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sql_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSqlServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlServerExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "version"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "administrator_login"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMSqlServer_basic(data acceptance.TestData) string {
	template := testAccAzureRMSqlServer_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_sql_server" "test" {
  name                = azurerm_sql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
