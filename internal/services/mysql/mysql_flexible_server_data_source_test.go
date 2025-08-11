// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MySQLFlexibleServerDataSource struct{}

func TestAccDataSourceMySqlFlexibleServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mysql_flexible_server", "test")
	r := MySQLFlexibleServerDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku_name").HasValue("B_Standard_B1ms"),
				check.That(data.ResourceName).Key("administrator_login").HasValue("_admin_Terraform_892123456789312"),
			),
		},
	})
}

func TestAccDataSourceMySqlFlexibleServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mysql_flexible_server", "test")
	r := MySQLFlexibleServerDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku_name").HasValue("MO_Standard_E2ds_v4"),
				check.That(data.ResourceName).Key("administrator_login").HasValue("adminTerraform"),
				check.That(data.ResourceName).Key("storage.0.size_gb").HasValue("20"),
				check.That(data.ResourceName).Key("version").HasValue("8.0.21"),
			),
		},
	})
}

func (MySQLFlexibleServerDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_mysql_flexible_server" "test" {
  name                = azurerm_mysql_flexible_server.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, MySqlFlexibleServerResource{}.basic(data))
}

func (MySQLFlexibleServerDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_mysql_flexible_server" "test" {
  name                = azurerm_mysql_flexible_server.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, MySqlFlexibleServerResource{}.complete(data))
}
