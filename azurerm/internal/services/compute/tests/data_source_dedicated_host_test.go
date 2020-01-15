package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAzureRMDedicatedHost_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dedicated_host", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedHost_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceDedicatedHost_basic(data acceptance.TestData) string {
	config := testAccAzureRMDedicatedHost_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_dedicated_host" "test" {
  name                = azurerm_dedicated_host.test.name
  resource_group_name = azurerm_dedicated_host.test.resource_group_name
  host_group_name     = azurerm_dedicated_host.test.host_group_name
}
`, config)
}
