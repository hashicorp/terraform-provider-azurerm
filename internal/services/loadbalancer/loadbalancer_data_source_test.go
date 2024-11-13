// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccAzureRMDataSourceLoadBalancer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_lb", "test")
	d := LoadBalancer{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.dataSourceBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("production"),
				check.That(data.ResourceName).Key("tags.Purpose").HasValue("AcceptanceTests"),
			),
		},
	})
}

func (r LoadBalancer) dataSourceBasic(data acceptance.TestData) string {
	resource := r.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_lb" "test" {
  name                = azurerm_lb.test.name
  resource_group_name = azurerm_lb.test.resource_group_name
}
`, resource)
}
