package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMVirtualHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_hub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVirtualHub_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "address_prefix"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_wan_id"),
				),
			},
		},
	})
}

func testAccDataSourceVirtualHub_basic(data acceptance.TestData) string {
	config := testAccAzureRMVirtualHub_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_hub" "test" {
  name                = azurerm_virtual_hub.test.name
  resource_group_name = azurerm_virtual_hub.test.resource_group_name
}
`, config)
}
