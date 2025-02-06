// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagerConnectivityConfigurationDataSource struct{}

func testAccNetworkManagerConnectivityConfigurationDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_manager_connectivity_configuration", "test")
	d := ManagerConnectivityConfigurationDataSource{}

	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("applies_to_group.#").HasValue("2"),
				check.That(data.ResourceName).Key("applies_to_group.0.global_mesh_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("applies_to_group.0.use_hub_gateway").HasValue("false"),
				check.That(data.ResourceName).Key("hub.#").HasValue("1"),
				check.That(data.ResourceName).Key("hub.0.resource_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("hub.0.resource_type").IsNotEmpty(),
			),
		},
	})
}

func (d ManagerConnectivityConfigurationDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_network_manager_connectivity_configuration" "test" {
  name               = azurerm_network_manager_connectivity_configuration.test.name
  network_manager_id = azurerm_network_manager_connectivity_configuration.test.network_manager_id
}
`, ManagerConnectivityConfigurationResource{}.complete(data))
}
