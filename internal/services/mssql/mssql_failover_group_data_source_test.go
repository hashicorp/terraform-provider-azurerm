// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MsSqlFailoverGroupDataSource struct{}

func TestAccMsSqlFailoverGroupDataSource_automaticFailover(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_failover_group", "test")
	r := MsSqlFailoverGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.automaticFailover(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("partner_server.#").HasValue("1"),
				check.That(data.ResourceName).Key("partner_server.0.id").Exists(),
				check.That(data.ResourceName).Key("read_write_endpoint_failover_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("read_write_endpoint_failover_policy.0.mode").HasValue("Automatic"),
				check.That(data.ResourceName).Key("read_write_endpoint_failover_policy.0.grace_minutes").HasValue("60"),
			),
		},
	})
}

func TestAccMsSqlFailoverGroupDataSource_automaticFailoverWithDatabases(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_failover_group", "test")
	r := MsSqlFailoverGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.automaticFailoverWithDatabases(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("partner_server.#").HasValue("1"),
				check.That(data.ResourceName).Key("partner_server.0.id").Exists(),
				check.That(data.ResourceName).Key("read_write_endpoint_failover_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("read_write_endpoint_failover_policy.0.mode").HasValue("Automatic"),
				check.That(data.ResourceName).Key("read_write_endpoint_failover_policy.0.grace_minutes").HasValue("80"),
				check.That(data.ResourceName).Key("databases.#").HasValue("1"),
				check.That(data.ResourceName).Key("databases.0").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("prod"),
				check.That(data.ResourceName).Key("tags.database").HasValue("test"),
			),
		},
	})
}

func (r MsSqlFailoverGroupDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test_primary" {
  name                         = "acctestmssql%[1]d-primary"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_server" "test_secondary" {
  name                         = "acctestmssql%[1]d-secondary"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = "%[3]s"
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "test" {
  name        = "acctestdb%[1]d"
  server_id   = azurerm_mssql_server.test_primary.id
  sku_name    = "S1"
  collation   = "SQL_Latin1_General_CP1_CI_AS"
  max_size_gb = "200"
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r MsSqlFailoverGroupDataSource) automaticFailover(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_failover_group" "test" {
  name      = "acctestsfg%[2]d"
  server_id = azurerm_mssql_server.test_primary.id

  partner_server {
    id = azurerm_mssql_server.test_secondary.id
  }

  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 60
  }
}

data "azurerm_mssql_failover_group" "test" {
  name      = azurerm_mssql_failover_group.test.name
  server_id = azurerm_mssql_server.test_primary.id
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlFailoverGroupDataSource) automaticFailoverWithDatabases(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_failover_group" "test" {
  name      = "acctestsfg%[2]d"
  server_id = azurerm_mssql_server.test_primary.id
  databases = [azurerm_mssql_database.test.id]

  partner_server {
    id = azurerm_mssql_server.test_secondary.id
  }

  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 80
  }

  tags = {
    environment = "prod"
    database    = "test"
  }
}

data "azurerm_mssql_failover_group" "test" {
  name      = azurerm_mssql_failover_group.test.name
  server_id = azurerm_mssql_server.test_primary.id
}
`, r.template(data), data.RandomInteger)
}
