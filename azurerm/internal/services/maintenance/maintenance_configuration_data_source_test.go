package maintenance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceMaintenanceConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_maintenance_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckMaintenanceConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMaintenanceConfiguration_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "scope", "Host"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "TesT"),
				),
			},
		},
	})
}

func testAccDataSourceMaintenanceConfiguration_complete(data acceptance.TestData) string {
	template := testAccMaintenanceConfiguration_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_maintenance_configuration" "test" {
  name                = azurerm_maintenance_configuration.test.name
  resource_group_name = azurerm_maintenance_configuration.test.resource_group_name
}
`, template)
}
