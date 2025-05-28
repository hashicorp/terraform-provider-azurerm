// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccDataSourceBackendAddressPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.dataSourceBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
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
