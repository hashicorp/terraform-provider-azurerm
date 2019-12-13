package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMPowerBIDedicatedCapacity_basic(t *testing.T) {
	dataSourceName := "data.azurerm_powerbidedicated_capacity.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePowerBIDedicatedCapacity_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "sku"),
				),
			},
		},
	})
}

func testAccDataSourcePowerBIDedicatedCapacity_basic(rInt int, location string) string {
	config := testAccAzureRMPowerBIDedicatedCapacity_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_powerbidedicated_capacity" "test" {
  name                = "${azurerm_powerbidedicated_capacity.test.name}"
  resource_group_name = "${azurerm_powerbidedicated_capacity.test.resource_group_name}"
}
`, config)
}
