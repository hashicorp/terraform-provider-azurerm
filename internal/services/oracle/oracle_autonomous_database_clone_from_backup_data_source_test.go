// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AutonomousDatabaseCloneFromBackupDataSource struct{}

func TestAccAutonomousDatabaseCloneFromBackupDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_autonomous_database_clone_from_backup", "test")
	r := AutonomousDatabaseCloneFromBackupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("database_version").Exists(),
				check.That(data.ResourceName).Key("compute_count").Exists(),
				check.That(data.ResourceName).Key("data_storage_size_in_tb").Exists(),
			),
		},
	})
}

func (AutonomousDatabaseCloneFromBackupDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_autonomous_database_clone_from_backup" "test" {
  name                = azurerm_oracle_autonomous_database_clone_from_backup.test.name
  resource_group_name = azurerm_oracle_autonomous_database_clone_from_backup.test.resource_group_name
}
`, AutonomousDatabaseCloneFromBackupResource{}.basic(data))
}
