package sql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type SqlServerDataSource struct{}

func TestAccDataSourceSqlServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sql_server", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: SqlServerDataSource{}.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("version").Exists(),
				check.That(data.ResourceName).Key("administrator_login").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func (d SqlServerDataSource) basic(data acceptance.TestData) string {
	template := SqlServerResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_sql_server" "test" {
  name                = azurerm_sql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
