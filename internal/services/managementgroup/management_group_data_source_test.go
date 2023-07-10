// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managementgroup_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagementGroupDataSource struct{}

func TestAccManagementGroupDataSource_basicByName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_management_group", "test")
	r := ManagementGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicByName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestmg-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("subscription_ids.#").HasValue("0"),
			),
		},
	})
}

func TestAccManagementGroupDataSource_basicByDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_management_group", "test")
	r := ManagementGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicByDisplayName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest Management Group %d", data.RandomInteger)),
				check.That(data.ResourceName).Key("subscription_ids.#").HasValue("0"),
			),
		},
	})
}

func TestAccManagementGroupDataSource_nestedManagmentGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_management_group", "test")
	r := ManagementGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.nestedManagementGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest Management Group %d", data.RandomInteger)),
				check.That(data.ResourceName).Key("management_group_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("all_management_group_ids.#").HasValue("2"),
			),
		},
	})
}

func (ManagementGroupDataSource) basicByName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  display_name = "acctestmg-%d"
}

data "azurerm_management_group" "test" {
  name = azurerm_management_group.test.name
}
`, data.RandomInteger)
}

func (ManagementGroupDataSource) basicByDisplayName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  display_name = "acctest Management Group %d"
}

data "azurerm_management_group" "test" {
  display_name = azurerm_management_group.test.display_name
}
`, data.RandomInteger)
}

func (ManagementGroupDataSource) nestedManagementGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  display_name = "acctest Management Group %[1]d"
}

resource "azurerm_management_group" "child" {
  display_name               = "acctest child Management Group %[1]d"
  parent_management_group_id = azurerm_management_group.test.id
}

resource "azurerm_management_group" "grand_child" {
  display_name               = "acctest grand child Management Group %[1]d"
  parent_management_group_id = azurerm_management_group.child.id
}

data "azurerm_management_group" "test" {
  name       = azurerm_management_group.test.name
  depends_on = [azurerm_management_group.grand_child]
}
`, data.RandomInteger)
}
