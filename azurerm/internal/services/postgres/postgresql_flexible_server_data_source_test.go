package postgres_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type PostgresqlFlexibleServerDataSource struct {
}

func TestAccDataSourcePostgresqlflexibleServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("administrator_login").Exists(),
				check.That(data.ResourceName).Key("sku_name").Exists(),
				check.That(data.ResourceName).Key("storage_mb").Exists(),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("backup_retention_days").Exists(),
				check.That(data.ResourceName).Key("cmk_enabled").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
	})
}

func (PostgresqlFlexibleServerDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_postgresql_flexible_server" "test" {
  name                = azurerm_postgresql_flexible_server.test.name
  resource_group_name = azurerm_postgresql_flexible_server.test.resource_group_name
}
`, PostgresqlFlexibleServerResource{}.basic(data))
}
