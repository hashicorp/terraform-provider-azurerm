package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAzureRMDataSourceLoadBalancerBackEndAddressPool_basic(t *testing.T) {
	dataSourceName := "data.azurerm_lb_backend_address_pool.test"
	ri := acctest.RandInt()
	location := testLocation()
	addressPoolName := fmt.Sprintf("%d-address-pool", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceLoadBalancerBackEndAddressPool_basic(ri, location, addressPoolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "backend_ip_configurations.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "load_balancing_rules.#", "0"),
				),
			},
		},
	})
}

func testAccAzureRMDataSourceLoadBalancerBackEndAddressPool_basic(rInt int, location, name string) string {
	resource := testAccAzureRMLoadBalancerBackEndAddressPool_basic(rInt, name, location)
	return fmt.Sprintf(`
%s

data "azurerm_lb_backend_address_pool" "test" {
  name            = "${azurerm_lb_backend_address_pool.test.name}"
  loadbalancer_id = "${azurerm_lb_backend_address_pool.test.loadbalancer_id}"
}
`, resource)
}
