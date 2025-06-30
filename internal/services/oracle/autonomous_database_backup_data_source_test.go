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
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("autonomous_database_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("backup_type").Exists(),
				check.That(data.ResourceName).Key("retention_period_in_days").Exists(),
				check.That(data.ResourceName).Key("autonomous_database_ocid").Exists(),
				check.That(data.ResourceName).Key("autonomous_database_backup_ocid").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
			),
		},
	})
}

func (r AutonomousDatabaseBackupDataSourceTest) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_autonomous_database_backup" "test" {
  name                     = azurerm_oracle_autonomous_database_backup.test.display_name
  resource_group_name      = azurerm_resource_group.test.name
  autonomous_database_name = azurerm_oracle_autonomous_database.test.name
}
`, AutonomousDatabaseBackupResource{}.complete(data))
}
