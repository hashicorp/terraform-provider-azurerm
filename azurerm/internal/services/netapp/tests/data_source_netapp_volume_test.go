package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMNetAppVolume_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_volume", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetAppVolume_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "volume_path"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "service_level"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_quota_in_gb"),
				),
			},
		},
	})
}

func testAccDataSourceNetAppVolume_basic(data acceptance.TestData) string {
	config := testAccAzureRMNetAppVolume_basic(data)
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
