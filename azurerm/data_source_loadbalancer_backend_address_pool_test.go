package azurerm

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAzureRMDataSourceLoadBalancerBackEndAddressPool_basic(t *testing.T) {
	dataSourceName := "data.azurerm_lb_backend_address_pool.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	addressPoolName := fmt.Sprintf("%d-address-pool", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceLoadBalancerBackEndAddressPool_basic(ri, location, addressPoolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
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
