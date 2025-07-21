// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
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
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("autonomous_database_id").Exists(),
			),
		},
	})
}

func (r AutonomousDatabaseBackupDataSourceTest) basic() string {
	return `

provider "azurerm" {
  features {}
}

data "azurerm_oracle_autonomous_database_backup" "test" {
  autonomous_database_id = "/subscriptions/4aa7be2d-ffd6-4657-828b-31ca25e39985/resourceGroups/dnsFarwoarder/providers/Oracle.Database/autonomousDatabases/DnsForwaderADBS"
}
`
}
