package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMNetAppVolume_basic(t *testing.T) {
	dataSourceName := "data.azurerm_netapp_volume.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetAppVolume_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "volume_path"),
					resource.TestCheckResourceAttrSet(dataSourceName, "service_level"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_quota_in_gb"),
				),
			},
		},
	})
}

func testAccDataSourceNetAppVolume_basic(rInt int, location string) string {
	config := testAccAzureRMNetAppVolume_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_netapp_volume" "test" {
  resource_group_name = "${azurerm_netapp_volume.test.resource_group_name}"
  account_name        = "${azurerm_netapp_volume.test.account_name}"
  pool_name           = "${azurerm_netapp_volume.test.pool_name}"
  name                = "${azurerm_netapp_volume.test.name}"
}
`, config)
}
