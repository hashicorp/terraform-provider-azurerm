package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMVirtualHub_basic(t *testing.T) {
	dataSourceName := "data.azurerm_virtual_hub.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVirtualHub_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "address_prefix"),
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_wan_id"),
				),
			},
		},
	})
}

func testAccDataSourceVirtualHub_basic(rInt int, location string) string {
	config := testAccAzureRMVirtualHub_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_hub" "test" {
  name                = azurerm_virtual_hub.test.name
  resource_group_name = azurerm_virtual_hub.test.resource_group_name
}
`, config)
}
