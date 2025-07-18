// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AutonomousDatabaseCloneDataSource struct{}

func TestAccAutonomousDatabaseCloneDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_autonomous_database_clone", "test")
	r := AutonomousDatabaseCloneDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("clone_type").HasValue("Full"),
				check.That(data.ResourceName).Key("lifecycle_state").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("db_version").Exists(),
				check.That(data.ResourceName).Key("compute_count").Exists(),
				check.That(data.ResourceName).Key("data_storage_size_in_tbs").Exists(),
			),
		},
	})
}

func TestAccAutonomousDatabaseCloneDataSource_metadataClone(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_oracle_autonomous_database_clone", "test")
	r := AutonomousDatabaseCloneDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.metadataClone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("clone_type").HasValue("Metadata"),
				check.That(data.ResourceName).Key("lifecycle_state").Exists(),
				check.That(data.ResourceName).Key("oci_url").Exists(),
			),
		},
	})
}

func (AutonomousDatabaseCloneDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_autonomous_database_clone" "test" {
  name                = azurerm_oracle_autonomous_database_clone.test.name
  resource_group_name = azurerm_oracle_autonomous_database_clone.test.resource_group_name
}
`, AutonomousDatabaseCloneResource{}.basic(data))
}

func (AutonomousDatabaseCloneDataSource) metadataClone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_autonomous_database_clone" "test" {
  name                = azurerm_oracle_autonomous_database_clone.test.name
  resource_group_name = azurerm_oracle_autonomous_database_clone.test.resource_group_name
}
`, AutonomousDatabaseCloneResource{}.metadataClone(data))
}
