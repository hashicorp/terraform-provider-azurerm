package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMNetAppSnapshot_basic(t *testing.T) {
	dataSourceName := "data.azurerm_netapp_snapshot.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetAppSnapshot_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceNetAppSnapshot_basic(rInt int, location string) string {
	config := testAccAzureRMNetAppSnapshot_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_netapp_snapshot" "test" {
  resource_group_name = "${azurerm_netapp_snapshot.test.resource_group_name}"
  account_name        = "${azurerm_netapp_snapshot.test.account_name}"
  pool_name           = "${azurerm_netapp_snapshot.test.pool_name}"
  volume_name         = "${azurerm_netapp_snapshot.test.volume_name}"
  name                = "${azurerm_netapp_snapshot.test.name}"
}
`, config)
}
