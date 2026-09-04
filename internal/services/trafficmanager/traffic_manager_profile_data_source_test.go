// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package trafficmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type TrafficManagerProfileDataSource struct{}

func TestAccTrafficManagerProfileDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_traffic_manager_profile", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: TrafficManagerProfileDataSource{}.template(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("traffic_routing_method").HasValue("Performance"),
			),
		},
	})
}

func TestAccTrafficManagerProfileDataSource_multiValue(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_traffic_manager_profile", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: TrafficManagerProfileDataSource{}.multiValue(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("max_return").HasValue("8"),
			),
		},
	})
}

func (d TrafficManagerProfileDataSource) template(data acceptance.TestData) string {
	template := TrafficManagerProfileResource{}.basic(data, "Performance")
	return fmt.Sprintf(`
%s

data "azurerm_traffic_manager_profile" "test" {
  name                = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func (d TrafficManagerProfileDataSource) multiValue(data acceptance.TestData) string {
	template := TrafficManagerProfileResource{}.multiValue(data)
	return fmt.Sprintf(`
%s

data "azurerm_traffic_manager_profile" "test" {
  name                = azurerm_traffic_manager_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
