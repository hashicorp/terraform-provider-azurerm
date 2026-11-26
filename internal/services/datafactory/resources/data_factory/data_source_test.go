// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DataFactoryDataSource struct{}

func TestAccDataFactoryDataSource_regressionTest(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_factory", "test")
	r := DataFactoryDataSource{}

	data.DataSourceRegressionTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
		},
	}, "")
}

func TestAccDataFactoryDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_factory", "test")
	r := DataFactoryDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
	})
}

func TestAccDataFactoryDataSource_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")
	r := DataFactoryDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("identity.#").Exists(),
				check.That(data.ResourceName).Key("identity.0.type").Exists(),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
	})
}

func TestAccDataFactoryDataSource_github(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_factory", "test")
	r := DataFactoryDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.github(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("github_configuration.#").HasValue("1"),
				check.That(data.ResourceName).Key("github_configuration.0.account_name").HasValue(fmt.Sprintf("acctestGH-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("github_configuration.0.last_commit_id").HasValue("0f9e8d7c6b5a4f3e2d1c0b9a8f7e6d5c4b3a2f1e"),
			),
		},
	})
}

func (DataFactoryDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_factory" "test" {
  name                = azurerm_data_factory.test.name
  resource_group_name = azurerm_data_factory.test.resource_group_name
}
`, DataFactoryResource{}.basic(data))
}

func (DataFactoryDataSource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_factory" "test" {
  name                = azurerm_data_factory.test.name
  resource_group_name = azurerm_data_factory.test.resource_group_name
}
`, DataFactoryResource{}.identity(data))
}

func (DataFactoryDataSource) github(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_factory" "test" {
  name                = azurerm_data_factory.test.name
  resource_group_name = azurerm_data_factory.test.resource_group_name
}
`, DataFactoryResource{}.github(data))
}
