package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMNetAppPool_basic(t *testing.T) {
	dataSourceName := "data.azurerm_netapp_pool.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetAppPool_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "account_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "service_level"),
					resource.TestCheckResourceAttrSet(dataSourceName, "size_in_tb"),
				),
			},
		},
	})
}

func testAccDataSourceNetAppPool_basic(rInt int, location string) string {
	config := testAccAzureRMNetAppPool_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_netapp_pool" "test" {
  resource_group_name = "${azurerm_netapp_pool.test.resource_group_name}"
  account_name        = "${azurerm_netapp_pool.test.account_name}"
  name                = "${azurerm_netapp_pool.test.name}"
}
`, config)
}
