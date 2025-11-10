// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AutonomousDatabaseBackupsDataSourceTest struct{}

func TestAccAutonomousDatabaseBackupsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_autonomous_database_backups", "test")
	r := AutonomousDatabaseBackupsDataSourceTest{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
	})
}

func (r AutonomousDatabaseBackupsDataSourceTest) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_autonomous_database_backups" "test" {
  autonomous_database_id = azurerm_oracle_autonomous_database.test.id
}
`, AutonomousDatabaseBackupResource{}.complete(data))
}
