package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAzureRMDataSourceLoadBalancerBackEndAddressPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_lb_backend_address_pool", "test")
	addressPoolName := fmt.Sprintf("%d-address-pool", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceLoadBalancerBackEndAddressPool_basic(data, addressPoolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
				),
			},
		},
	})
}

func TestAccAzureRMDataSourceLoadBalancerBackEndAddressPool_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_lb_backend_address_pool", "test")
	addressPoolName := fmt.Sprintf("%d-address-pool", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceLoadBalancerBackEndAddressPool_standard(data, addressPoolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_address.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_address.0.name", "addr-1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address.0.virtual_network_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_address.0.ip_address", "10.0.1.4"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_address.1.name", "addr-2"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address.1.virtual_network_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_address.1.ip_address", "10.0.1.5"),
				),
			},
		},
	})
}

func testAccAzureRMDataSourceLoadBalancerBackEndAddressPool_basic(data acceptance.TestData, name string) string {
	resource := testAccAzureRMLoadBalancerBackEndAddressPool_basic(data, name)
	return fmt.Sprintf(`
%s

data "azurerm_lb_backend_address_pool" "test" {
  name            = azurerm_lb_backend_address_pool.test.name
  loadbalancer_id = azurerm_lb_backend_address_pool.test.loadbalancer_id
}
`, resource)
}

func testAccAzureRMDataSourceLoadBalancerBackEndAddressPool_standard(data acceptance.TestData, name string) string {
	resource := testAccAzureRMLoadBalancerBackEndAddressPool_standard(data, name, "Standard", true)
	return fmt.Sprintf(`
%s

data "azurerm_lb_backend_address_pool" "test" {
  name            = azurerm_lb_backend_address_pool.test.name
  loadbalancer_id = azurerm_lb_backend_address_pool.test.loadbalancer_id
}
`, resource)
}
