package postgres_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type PostgreSQLServerDataSource struct {
}

func TestAccDataSourcePostgreSqlServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_postgresql_server", "test")
	r := PostgreSQLServerDataSource{}
	version := "9.5"

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data, version),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("administrator_login").Exists(),
				check.That(data.ResourceName).Key("sku_name").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func (r PostgreSQLServerDataSource) basic(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

data "azurerm_postgresql_server" "test" {
  name                = azurerm_postgresql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, PostgreSQLServerResource{}.basic(data, version))
}
