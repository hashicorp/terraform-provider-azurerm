// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type KubernetesFleetManagerDataSource struct{}

func TestAccKubernetesFleetManagerDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_fleet_manager", "test")
	d := KubernetesFleetManagerDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("hub_profile.#").HasValue("1"),
				check.That(data.ResourceName).Key("hub_profile.0.agent_profile.#").HasValue("1"),
				check.That(data.ResourceName).Key("hub_profile.0.agent_profile.0.subnet_id").HasValue(""),
				check.That(data.ResourceName).Key("hub_profile.0.agent_profile.0.virtual_machine_size").HasValue("Standard_DS2_v2"),
				check.That(data.ResourceName).Key("hub_profile.0.api_server_access_profile.#").HasValue("1"),
				check.That(data.ResourceName).Key("hub_profile.0.api_server_access_profile.0.enable_private_cluster").HasValue("false"),
				check.That(data.ResourceName).Key("hub_profile.0.dns_prefix").HasValue(fmt.Sprintf("val-%s", data.RandomString)),
				check.That(data.ResourceName).Key("hub_profile.0.fqdn").Exists(),
				check.That(data.ResourceName).Key("hub_profile.0.kubernetes_version").Exists(),
				check.That(data.ResourceName).Key("hub_profile.0.portal_fqdn").Exists(),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("terraform-acctests"),
				check.That(data.ResourceName).Key("tags.some_key").HasValue("some-value"),
			),
		},
	})
}

func (KubernetesFleetManagerDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_fleet_manager" "test" {
  name                = azurerm_kubernetes_fleet_manager.test.name
  resource_group_name = azurerm_kubernetes_fleet_manager.test.resource_group_name
}
`, KubernetesFleetManagerTestResource{}.complete(data))
}
