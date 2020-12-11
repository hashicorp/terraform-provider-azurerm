package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAzureRMDataSourceLoadBalancerBackEndAddressPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_lb_backend_address_pool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceLoadBalancerBackEndAddressPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
				),
			},
		},
	})
}

func testAccAzureRMDataSourceLoadBalancerBackEndAddressPool_basic(data acceptance.TestData) string {
	resource := testAccAzureRMLoadBalancerBackEndAddressPool_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_lb_backend_address_pool" "test" {
  name            = azurerm_lb_backend_address_pool.test.name
  loadbalancer_id = azurerm_lb_backend_address_pool.test.loadbalancer_id
}
`, resource)
}
