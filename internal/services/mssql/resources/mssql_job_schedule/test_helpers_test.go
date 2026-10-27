// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_schedule_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type MsSqlJobResource struct{}

func (r MsSqlJobResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_job" "test" {
  name         = "acctest-job-%d"
  job_agent_id = azurerm_mssql_job_agent.test.id
}
`, r.template(data), data.RandomInteger)
}

func (MsSqlJobResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctest-server-%[1]d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[1]d"
  server_id = azurerm_mssql_server.test.id
  collation = "SQL_Latin1_General_CP1_CI_AS"
  sku_name  = "S1"
}

resource "azurerm_mssql_job_agent" "test" {
  name        = "acctest-job-agent-%[1]d"
  location    = azurerm_resource_group.test.location
  database_id = azurerm_mssql_database.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}
