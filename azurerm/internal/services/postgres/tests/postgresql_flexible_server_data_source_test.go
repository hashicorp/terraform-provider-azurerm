package tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"testing"
)

func TestAccDataSourceAzureRMpostgresqlflexibleServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_postgresql_flexible_server", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMpostgresqlflexibleServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcepostgresqlflexibleServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMpostgresqlflexibleServerExists(data.ResourceName),
				),
			},
		},
	})
}

func testAccDataSourcepostgresqlflexibleServer_basic(data acceptance.TestData) string {
	config := testAccAzureRMpostgresqlflexibleServer_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_postgresql_flexible_server" "test" {
  name                = azurerm_postgresql_flexible_server.test.name
  resource_group_name = azurerm_postgresql_flexible_server.test.resource_group_name
}
`, config)
}
