// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AutonomousDatabaseBackupDataSourceTest struct{}

func TestAccAutonomousDatabaseBackupDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_autonomous_database_backup", "test")
	r := AutonomousDatabaseBackupDataSourceTest{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("autonomous_database_id").Exists(),
			),
		},
	})
}

func (r AutonomousDatabaseBackupDataSourceTest) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_autonomous_database_backup" "test" {
  autonomous_database_id = azurerm_oracle_autonomous_database.test.id
}
`, AutonomousDatabaseBackupResource{}.complete(data))
}
