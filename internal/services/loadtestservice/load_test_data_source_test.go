// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadtestservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type LoadTestTestDataSource struct{}

func TestAccLoadTestDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_load_test", "test")
	d := LoadTestTestDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("description").HasValue("Description for the Load Test"),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned, UserAssigned"),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("terraform-acctests"),
				check.That(data.ResourceName).Key("tags.some_key").HasValue("some-value"),
			),
		},
	})
}

func TestAccLoadTestDataSource_encryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_load_test", "test")
	d := LoadTestTestDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.encryption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("encryption.0.identity.0.type").HasValue("UserAssigned"),
			),
		},
	})
}

func (d LoadTestTestDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data azurerm_load_test test {
  name                = azurerm_load_test.test.name
  resource_group_name = azurerm_load_test.test.resource_group_name
}
`, LoadTestTestResource{}.complete(data))
}

func (d LoadTestTestDataSource) encryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data azurerm_load_test test {
  name                = azurerm_load_test.test.name
  resource_group_name = azurerm_load_test.test.resource_group_name
}
`, LoadTestTestResource{}.encryption(data))
}
