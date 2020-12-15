package loadbalancer_test

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAzureRMDataSourceLoadBalancerBackEndAddressPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_lb_backend_address_pool", "test")
	d := LoadBalancerBackendAddressPool{}


	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: d.dataSourceBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
	})
}

func (d LoadBalancerBackendAddressPool) dataSourceBasic(data acceptance.TestData) string {
	resource := d.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_lb_backend_address_pool" "test" {
  name            = azurerm_lb_backend_address_pool.test.name
  loadbalancer_id = azurerm_lb_backend_address_pool.test.loadbalancer_id
}
`, resource)
}
