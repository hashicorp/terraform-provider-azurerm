package mysql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type MySQLServerDataSource struct {
}

func TestAccDataSourceMySQLServerDataSourceMySQLServer_basicFiveSix(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mysql_server", "test")
	r := MySQLServerDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data, "5.6"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
				check.That(data.ResourceName).Key("administrator_login").HasValue("acctestun"),
				check.That(data.ResourceName).Key("auto_grow_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("ssl_minimal_tls_version_enforced").HasValue("TLS1_1"),
				check.That(data.ResourceName).Key("storage_mb").HasValue("51200"),
				check.That(data.ResourceName).Key("version").HasValue("5.6"),
			),
		},
	})
}

func TestAccDataSourceMySQLServerDataSourceMySQLServer_basicFiveSixWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mysql_server", "test")
	r := MySQLServerDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicWithIdentity(data, "5.6"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
				check.That(data.ResourceName).Key("administrator_login").HasValue("acctestun"),
				check.That(data.ResourceName).Key("auto_grow_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("ssl_minimal_tls_version_enforced").HasValue("TLS1_1"),
				check.That(data.ResourceName).Key("storage_mb").HasValue("51200"),
				check.That(data.ResourceName).Key("version").HasValue("5.6"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceMySQLServerDataSourceMySQLServer_basicFiveSeven(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mysql_server", "test")
	r := MySQLServerDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data, "5.7"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
				check.That(data.ResourceName).Key("administrator_login").HasValue("acctestun"),
				check.That(data.ResourceName).Key("auto_grow_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("ssl_minimal_tls_version_enforced").HasValue("TLS1_1"),
				check.That(data.ResourceName).Key("storage_mb").HasValue("51200"),
				check.That(data.ResourceName).Key("version").HasValue("5.7"),
			),
		},
	})
}

func TestAccDataSourceMySQLServerDataSourceMySQLServer_basicEightZero(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mysql_server", "test")
	r := MySQLServerDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data, "8.0"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
				check.That(data.ResourceName).Key("administrator_login").HasValue("acctestun"),
				check.That(data.ResourceName).Key("auto_grow_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("ssl_minimal_tls_version_enforced").HasValue("TLS1_1"),
				check.That(data.ResourceName).Key("storage_mb").HasValue("51200"),
				check.That(data.ResourceName).Key("version").HasValue("8.0"),
			),
		},
	})
}

func TestAccDataSourceMySQLServerDataSourceMySQLServer_autogrowOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mysql_server", "test")
	r := MySQLServerDataSource{}
	mysqlVersion := "5.7"

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.autogrow(data, mysqlVersion),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
				check.That(data.ResourceName).Key("administrator_login").HasValue("acctestun"),
				check.That(data.ResourceName).Key("auto_grow_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("storage_mb").HasValue("51200"),
				check.That(data.ResourceName).Key("version").HasValue("5.7"),
			),
		},
	})
}

func TestAccDataSourceMySQLServerDataSourceMySQLServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mysql_server", "test")
	r := MySQLServerDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.complete(data, "8.0"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
				check.That(data.ResourceName).Key("administrator_login").HasValue("acctestun"),
				check.That(data.ResourceName).Key("auto_grow_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("ssl_minimal_tls_version_enforced").HasValue("TLS1_2"),
				check.That(data.ResourceName).Key("storage_mb").HasValue("51200"),
				check.That(data.ResourceName).Key("version").HasValue("8.0"),
				check.That(data.ResourceName).Key("threat_detection_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.enabled").HasValue("true"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.email_account_admins").HasValue("true"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.retention_days").HasValue("7"),
			),
		},
	})
}

func (MySQLServerDataSource) basic(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

data "azurerm_mysql_server" "test" {
  name                = azurerm_mysql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, MySQLServerResource{}.basic(data, version))
}

func (MySQLServerDataSource) basicWithIdentity(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

data "azurerm_mysql_server" "test" {
  name                = azurerm_mysql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, MySQLServerResource{}.basicWithIdentity(data, version))
}

func (MySQLServerDataSource) autogrow(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

data "azurerm_mysql_server" "test" {
  name                = azurerm_mysql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, MySQLServerResource{}.autogrow(data, version))
}

func (MySQLServerDataSource) complete(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

data "azurerm_mysql_server" "test" {
  name                = azurerm_mysql_server.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, MySQLServerResource{}.complete(data, version))
}
