package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMPrivateLinkService_complete(t *testing.T) {
	dataSourceName := "data.azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivateLinkService_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "primary_nat_ip_configuration.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "secondary_nat_ip_configuration.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "primary_nat_ip_configuration.0.private_ip_address", "10.5.1.17"),
					resource.TestCheckResourceAttr(dataSourceName, "primary_nat_ip_configuration.0.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(dataSourceName, "secondary_nat_ip_configuration.0.private_ip_address", "10.5.1.18"),
					resource.TestCheckResourceAttr(dataSourceName, "secondary_nat_ip_configuration.0.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(dataSourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourcePrivateLinkService_complete(rInt int, location string) string {
	config := testAccAzureRMPrivateLinkService_complete(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_private_link_service" "test" {
  name                = azurerm_private_link_service.test.name
  resource_group_name = azurerm_private_link_service.test.resource_group_name
}
`, config)
}
