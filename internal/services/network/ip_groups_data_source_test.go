// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type IPGroupsDataSource struct{}

func TestAccDataSourceIPGroups_noResults(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_ip_groups", "test")
	r := IPGroupsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.noResults(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("ids.#").HasValue("0"),
				check.That(data.ResourceName).Key("names.#").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceIPGroups_single(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_ip_groups", "test")
	r := IPGroupsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.single(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("names.#").HasValue("1"),
			),
		},
	})
}

func TestAccDataSourceIPGroups_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_ip_groups", "test")
	r := IPGroupsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("ids.#").HasValue("3"),
				check.That(data.ResourceName).Key("names.#").HasValue("3"),
			),
		},
	})
}

// Find IP group which doesn't exist
func (IPGroupsDataSource) noResults(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_ip_groups" "test" {
  name                = "doesNotExist"
  resource_group_name = azurerm_resource_group.test.name
  depends_on = [
    azurerm_ip_group.test,
  ]
}
`, IPGroupResource{}.basic(data))
}

// Find single IP group
func (IPGroupsDataSource) single(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_ip_groups" "test" {
  name                = "acceptanceTestIpGroup1"
  resource_group_name = azurerm_resource_group.test.name
  depends_on = [
    azurerm_ip_group.test,
  ]
}
`, IPGroupResource{}.basic(data))
}

// Find multiple IP Groups, filtered by name substring
func (IPGroupsDataSource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_ip_groups" "test" {
  name                = "acceptanceTestIpGroup"
  resource_group_name = azurerm_resource_group.test.name
  depends_on = [
    azurerm_ip_group.test,
    azurerm_ip_group.test2,
    azurerm_ip_group.test3,
  ]
}
`, IPGroupResource{}.complete(data))
}
