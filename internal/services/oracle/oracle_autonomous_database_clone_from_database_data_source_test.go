// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AutonomousDatabaseCloneFromDatabaseDataSource struct{}

func TestAccAutonomousDatabaseCloneFromDatabaseDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_autonomous_database_clone_from_database", "test")
	r := AutonomousDatabaseCloneFromDatabaseDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("lifecycle_state").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("compute_count").Exists(),
				check.That(data.ResourceName).Key("data_storage_size_in_tb").Exists(),
			),
		},
	})
}

func TestAccAutonomousDatabaseCloneFromDatabaseDataSource_metadataClone(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_autonomous_database_clone_from_database", "test")
	r := AutonomousDatabaseCloneFromDatabaseDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.metadataClone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("lifecycle_state").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("compute_count").Exists(),
				check.That(data.ResourceName).Key("data_storage_size_in_tb").Exists(),
			),
		},
	})
}

func (AutonomousDatabaseCloneFromDatabaseDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_autonomous_database_clone_from_database" "test" {
  name                = azurerm_oracle_autonomous_database_clone_from_database.test.name
  resource_group_name = azurerm_oracle_autonomous_database_clone_from_database.test.resource_group_name
}
`, AutonomousDatabaseCloneFromDatabaseResource{}.basic(data))
}

func (AutonomousDatabaseCloneFromDatabaseDataSource) metadataClone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_autonomous_database_clone_from_database" "test" {
  name                = azurerm_oracle_autonomous_database_clone_from_database.test.name
  resource_group_name = azurerm_oracle_autonomous_database_clone_from_database.test.resource_group_name
}
`, AutonomousDatabaseCloneFromDatabaseResource{}.metadataClone(data))
}
