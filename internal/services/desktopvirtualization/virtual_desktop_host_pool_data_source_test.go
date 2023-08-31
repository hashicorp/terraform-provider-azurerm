// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type VirtualDesktopHostPoolDataSource struct{}

func TestAccDataShareAccountDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_desktop_host_pool", "test")
	r := VirtualDesktopHostPoolDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("friendly_name").HasValue("A Friendly Name!"),
				check.That(data.ResourceName).Key("validate_environment").HasValue("true"),
				check.That(data.ResourceName).Key("load_balancer_type").HasValue("BreadthFirst"),
				check.That(data.ResourceName).Key("maximum_sessions_allowed").HasValue("100"),
			),
		},
	})
}

func (VirtualDesktopHostPoolDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_virtual_desktop_host_pool" "test" {
  name                = azurerm_virtual_desktop_host_pool.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, VirtualDesktopHostPoolResource{}.complete(data))
}
