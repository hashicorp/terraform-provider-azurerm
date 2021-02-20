package postgres_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"testing"
)

type PostgresqlFlexibleServerDataSource struct {
}

func TestAccDataSourcePostgresqlflexibleServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("administrator_login").Exists(),
				check.That(data.ResourceName).Key("sku.0.name").Exists(),
				check.That(data.ResourceName).Key("sku.0.tier").Exists(),
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
