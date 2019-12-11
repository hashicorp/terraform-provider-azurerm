package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMPrivateLinkService_complete(t *testing.T) {
	dataSourceName := "data.azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivateLinkService_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "nat_ip_configuration.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "nat_ip_configuration.0.private_ip_address", "10.5.1.40"),
					resource.TestCheckResourceAttr(dataSourceName, "nat_ip_configuration.0.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(dataSourceName, "nat_ip_configuration.1.private_ip_address", "10.5.1.41"),
					resource.TestCheckResourceAttr(dataSourceName, "nat_ip_configuration.1.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(dataSourceName, "auto_approval_subscription_ids.0", subscriptionId),
					resource.TestCheckResourceAttr(dataSourceName, "visibility_subscription_ids.0", subscriptionId),
					resource.TestCheckResourceAttr(dataSourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "load_balancer_frontend_ip_configuration_ids.0"),
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
