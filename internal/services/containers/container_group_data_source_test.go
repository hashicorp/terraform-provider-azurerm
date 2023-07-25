// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ContainerGroupDataSource struct{}

func TestAccDataSourceContainerGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_group", "test")
	r := ContainerGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
			),
		},
	})
}

func TestAccDataSourceContainerGroup_UAMIWithSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_group", "test")
	r := ContainerGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.uamiWithSubnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("subnet_ids.#").HasValue("1"),
			),
		},
	})
}

func (ContainerGroupDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_group" "test" {
  name                = azurerm_container_group.test.name
  resource_group_name = azurerm_container_group.test.resource_group_name
}
`, ContainerGroupResource{}.linuxComplete(data))
}

func (ContainerGroupDataSource) uamiWithSubnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_group" "test" {
  name                = azurerm_container_group.test.name
  resource_group_name = azurerm_container_group.test.resource_group_name
}
`, ContainerGroupResource{}.UserAssignedIdentityWithVirtualNetwork(data))
}
