package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMNetAppAccount_basic(t *testing.T) {
	dataSourceName := "data.azurerm_netapp_account.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetAppAccount_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceNetAppAccount_basic(rInt int, location string) string {
	config := testAccAzureRMNetAppAccount_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_netapp_account" "test" {
  resource_group_name = "${azurerm_netapp_account.test.resource_group_name}"
  name                = "${azurerm_netapp_account.test.name}"
}
`, config)
}
