package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMPrivateLinkService_basic(t *testing.T) {
	dataSourceName := "data.azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivateLinkService_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "fqdns.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "fqdns.0", "testFqdns"),
					resource.TestCheckResourceAttr(dataSourceName, "ip_configurations.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "ip_configurations.0.private_ip_address", "10.5.1.17"),
					resource.TestCheckResourceAttr(dataSourceName, "ip_configurations.0.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(dataSourceName, "ip_configurations.0.private_ip_allocation_method", "Static"),
					resource.TestCheckResourceAttr(dataSourceName, "load_balancer_frontend_ip_configurations.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourcePrivateLinkService_basic(rInt int, location string) string {
	config := testAccAzureRMPrivateLinkService_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_private_link_service" "test" {
  resource_group_name = "${azurerm_private_link_service.test.resource_group_name}"
  name                = "${azurerm_private_link_service.test.name}"
}
`, config)
}
