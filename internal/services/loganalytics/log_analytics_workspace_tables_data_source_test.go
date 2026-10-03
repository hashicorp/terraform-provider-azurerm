// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type LogAnalyticsWorkspaceTablesDataSource struct{}

func TestAccDataSourceLogAnalyticsWorkspaceTables_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_log_analytics_workspace_tables", "test")
	r := LogAnalyticsWorkspaceTablesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttrSet(data.ResourceName, "tables.#"),
				acceptance.TestCheckResourceAttrSet(data.ResourceName, "names.#"),

				check.That(data.ResourceName).Key("tables.0.name").Exists(),
				check.That(data.ResourceName).Key("tables.0.plan").Exists(),
				check.That(data.ResourceName).Key("tables.0.retention_in_days").Exists(),
				check.That(data.ResourceName).Key("tables.0.total_retention_in_days").Exists(),

				check.That(data.ResourceName).Key("names.0").Exists(),
			),
		},
	})
}

func (LogAnalyticsWorkspaceTablesDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  retention_in_days   = 30
}

data "azurerm_log_analytics_workspace_tables" "test" {
  workspace_id = azurerm_log_analytics_workspace.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
