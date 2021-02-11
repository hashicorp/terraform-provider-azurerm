package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

func TestAccAzureRMDataSourceLoadBalancerBackEndAddressPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.dataSourceBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
	})
}

func (r LoadBalancerBackendAddressPool) dataSourceBasic(data acceptance.TestData) string {
	resource := r.basicSkuBasic(data)
	return fmt.Sprintf(`
%s

data "azurerm_lb_backend_address_pool" "test" {
  name            = azurerm_lb_backend_address_pool.test.name
  loadbalancer_id = azurerm_lb_backend_address_pool.test.loadbalancer_id
}
`, resource)
}
