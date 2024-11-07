// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devtestlabs_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DevTestLabDataSource struct{}

func TestAccDevTestLabDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dev_test_lab", "test")
	r := DevTestLabDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("storage_type").HasValue("Premium"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDevTestLabDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dev_test_lab", "test")
	r := DevTestLabDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("storage_type").HasValue("Premium"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Hello").HasValue("World"),
			),
		},
	})
}

func (DevTestLabDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_dev_test_lab" "test" {
  name                = azurerm_dev_test_lab.test.name
  resource_group_name = azurerm_dev_test_lab.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (DevTestLabDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    Hello = "World"
  }
}

data "azurerm_dev_test_lab" "test" {
  name                = azurerm_dev_test_lab.test.name
  resource_group_name = azurerm_dev_test_lab.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
