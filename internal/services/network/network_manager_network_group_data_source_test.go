// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagerNetworkGroupDataSource struct{}

func testAccNetworkManagerNetworkGroupDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_manager_network_group", "test")
	d := ManagerNetworkGroupDataSource{}
	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("description").HasValue("test complete"),
			),
		},
	})
}

func (d ManagerNetworkGroupDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_network_manager_network_group" "test" {
  name               = azurerm_network_manager_network_group.test.name
  network_manager_id = azurerm_network_manager_network_group.test.network_manager_id
}
`, ManagerNetworkGroupResource{}.complete(data))
}
