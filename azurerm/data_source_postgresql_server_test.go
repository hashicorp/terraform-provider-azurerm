package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMPPostgreSqlServer_basic(t *testing.T) {
	dataSourceName := "data.azurerm_postgresql_server.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	version := "9.5"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPostgreSQLServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMPostgreSqlServer_basic(ri, location, version),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPostgreSQLServerExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttrSet(dataSourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "administrator_login"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMPostgreSqlServer_basic(rInt int, location string, version string) string {
	template := testAccAzureRMPostgreSQLServer_basic(rInt, location, version)
	return fmt.Sprintf(`
%s

data "azurerm_postgresql_server" "test" {
  name                = "${azurerm_postgresql_server.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}
