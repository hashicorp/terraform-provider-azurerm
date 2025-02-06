// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagerDataSource struct{}

func testAccNetworkManagerDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_manager", "test")
	d := ManagerDataSource{}
	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").IsNotEmpty(),
				check.That(data.ResourceName).Key("description").IsNotEmpty(),
				check.That(data.ResourceName).Key("scope_accesses.#").HasValue("2"),
				check.That(data.ResourceName).Key("scope_accesses.0").HasValue("Connectivity"),
				check.That(data.ResourceName).Key("scope.#").HasValue("1"),
				check.That(data.ResourceName).Key("scope.0.subscription_ids.#").HasValue("1"),
			),
		},
	})
}

func (d ManagerDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_network_manager" "test" {
  name                = azurerm_network_manager.test.name
  resource_group_name = azurerm_network_manager.test.resource_group_name
}
`, ManagerResource{}.complete(data))
}
