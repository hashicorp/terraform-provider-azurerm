package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMDataMigrationService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_migration_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataMigrationService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Standard_1vCores"),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "Cloud"),
				),
			},
		},
	})
}

func testAccDataSourceDataMigrationService_basic(data acceptance.TestData) string {
	config := testAccAzureRMDataMigrationService_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_data_migration_service" "test" {
  resource_group_name   = "${azurerm_data_migration_service.test.resource_group_name}"
  name                  = "${azurerm_data_migration_service.test.name}"
}
`, config)
}
