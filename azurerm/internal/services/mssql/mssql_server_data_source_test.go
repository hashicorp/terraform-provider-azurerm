package mssql_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type MsSqlServerDataSource struct{}

func TestAccDataSourceMsSqlServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_server", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: MsSqlServerDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestsqlserver%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
			),
		},
	})
}

func TestAccDataSourceMsSqlServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_server", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: MsSqlServerDataSource{}.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestsqlserver%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("administrator_login").Exists(),
				check.That(data.ResourceName).Key("fully_qualified_domain_name").Exists(),
				check.That(data.ResourceName).Key("tags.%").Exists(),
			),
		},
	})
}

func (MsSqlServerDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_mssql_server" "test" {
  name                = azurerm_mssql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
}

`, MsSqlServerResource{}.basic(data))
}

func (MsSqlServerDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_mssql_server" "test" {
  name                = azurerm_mssql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
}

`, MsSqlServerResource{}.complete(data))
}
