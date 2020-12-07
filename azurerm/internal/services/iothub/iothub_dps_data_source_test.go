package iothub

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMIotHubDPS_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_iothub_dps", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDPSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMIotHubDPS_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "allocation_policy"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "device_provisioning_host_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "id_scope"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "service_operations_host_name"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMIotHubDPS_basic(data acceptance.TestData) string {
	template := testAccAzureRMIotHubDPS_basic(data)

	return fmt.Sprintf(`
%s

data "azurerm_iothub_dps" "test" {
  name                = azurerm_iothub_dps.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
