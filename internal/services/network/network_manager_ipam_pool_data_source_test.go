// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagerIpamPoolDataSource struct{}

func testAccNetworkManagerIpamPoolDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_manager_ipam_pool", "test")
	d := ManagerIpamPoolDataSource{}
	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("address_prefixes.#").HasValue("2"),
				check.That(data.ResourceName).Key("description").HasValue("This is a test IPAM pool"),
				check.That(data.ResourceName).Key("display_name").HasValue("ipampool1"),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("parent_pool_name").IsNotEmpty(),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
	})
}

func (d ManagerIpamPoolDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_network_manager_ipam_pool" "test" {
  name               = azurerm_network_manager_ipam_pool.test.name
  network_manager_id = azurerm_network_manager_ipam_pool.test.network_manager_id
}
`, ManagerIpamPoolResource{}.complete(data))
}
