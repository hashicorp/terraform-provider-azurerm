package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMPrivateLinkService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_link_service", "test")
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivateLinkService_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.0.private_ip_address", "10.5.1.40"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.0.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.1.private_ip_address", "10.5.1.41"),
					resource.TestCheckResourceAttr(data.ResourceName, "nat_ip_configuration.1.private_ip_address_version", "IPv4"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_approval_subscription_ids.0", subscriptionId),
					resource.TestCheckResourceAttr(data.ResourceName, "visibility_subscription_ids.0", subscriptionId),
					resource.TestCheckResourceAttr(data.ResourceName, "load_balancer_frontend_ip_configuration_ids.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "load_balancer_frontend_ip_configuration_ids.0"),
				),
			},
		},
	})
}

func testAccDataSourcePrivateLinkService_complete(data acceptance.TestData) string {
	config := testAccAzureRMPrivateLinkService_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_private_link_service" "test" {
  name                = azurerm_private_link_service.test.name
  resource_group_name = azurerm_private_link_service.test.resource_group_name
}
`, config)
}
