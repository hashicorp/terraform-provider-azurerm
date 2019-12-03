package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMDedicatedHostGroup_basic(t *testing.T) {
	dataSourceName := "data.azurerm_dedicated_host_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	rName := acctest.RandStringFromCharSet(4, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedHostGroup_basic(ri, location, rName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostGroupExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "zones.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "zones.0", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "platform_fault_domain_count", "2"),
				),
			},
		},
	})
}

func testAccDataSourceDedicatedHostGroup_basic(rInt int, location string, rName string) string {
	config := testAccAzureRMDedicatedHostGroup_complete(rInt, location, rName)
	return fmt.Sprintf(`
%s

data "azurerm_dedicated_host_group" "test" {
  resource_group_name 	= "${azurerm_dedicated_host_group.test.resource_group_name}"
  name           		= "${azurerm_dedicated_host_group.test.name}"
}
`, config)
}
