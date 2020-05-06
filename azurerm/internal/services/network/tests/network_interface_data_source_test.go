package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceArmNetworkInterface_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_interface", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterface_static(data),
			},
			{
				Config: testAccDataSourceNetworkInterface_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "private_ip_address", "10.0.2.15"),
				),
			},
		},
	})
}

func testAccDataSourceNetworkInterface_basic(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterface_static(data)
	return fmt.Sprintf(`
%s

data "azurerm_network_interface" "test" {
  name                = azurerm_network_interface.test.name
  resource_group_name = azurerm_network_interface.test.resource_group_name
}
`, template)
}
