// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PostgresqlFlexibleServerActiveDirectoryAdministratorDataSource struct{}

func TestAccPostgresqlFlexibleServerActiveDirectoryAdministratorDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_postgresql_flexible_server_active_directory_administrator", "test")
	r := PostgresqlFlexibleServerActiveDirectoryAdministratorDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("server_id").Exists(),
				check.That(data.ResourceName).Key("object_id").Exists(),
				check.That(data.ResourceName).Key("principal_name").Exists(),
				check.That(data.ResourceName).Key("principal_type").Exists(),
				check.That(data.ResourceName).Key("tenant_id").Exists(),
			),
		},
	})
}

func (PostgresqlFlexibleServerActiveDirectoryAdministratorDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_postgresql_flexible_server_active_directory_administrator" "test" {
  server_id = azurerm_postgresql_flexible_server.test.id
  object_id = azurerm_postgresql_flexible_server_active_directory_administrator.test.object_id
}
`, PostgresqlFlexibleServerAdministratorResource{}.basic(data))
}
