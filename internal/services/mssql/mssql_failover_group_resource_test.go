// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/failovergroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MsSqlFailoverGroupResource struct{}

func TestAccMsSqlFailoverGroup_automaticFailover(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_failover_group", "test")
	r := MsSqlFailoverGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.automaticFailover(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlFailoverGroup_automaticFailoverWithDatabases(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_failover_group", "test")
	r := MsSqlFailoverGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.automaticFailoverWithDatabases(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlFailoverGroup_manualFailover(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_failover_group", "test")
	r := MsSqlFailoverGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manualFailover(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlFailoverGroup_manualFailoverWithDatabases(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_failover_group", "test")
	r := MsSqlFailoverGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.manualFailoverWithDatabases(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlFailoverGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_failover_group", "test")
	r := MsSqlFailoverGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.automaticFailover(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.manualFailoverWithDatabases(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.automaticFailover(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlFailoverGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_failover_group", "test")
	r := MsSqlFailoverGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.automaticFailoverWithDatabases(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r MsSqlFailoverGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := failovergroups.ParseFailoverGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.FailoverGroupsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(true), nil
}

func (r MsSqlFailoverGroupResource) template(data acceptance.TestData) string {
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

func (r MsSqlFailoverGroupResource) automaticFailover(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r MsSqlFailoverGroupResource) manualFailover(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_failover_group" "test" {
  name      = "acctestsfg%[2]d"
  server_id = azurerm_mssql_server.test_primary.id

  partner_server {
    id = azurerm_mssql_server.test_secondary.id
  }

  read_write_endpoint_failover_policy {
    mode = "Manual"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlFailoverGroupResource) automaticFailoverWithDatabases(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r MsSqlFailoverGroupResource) manualFailoverWithDatabases(data acceptance.TestData) string {
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
    mode = "Manual"
  }

  tags = {
    environment = "staging"
    database    = "test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlFailoverGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_failover_group" "import" {
  name      = azurerm_mssql_failover_group.test.name
  server_id = azurerm_mssql_failover_group.test.server_id
  databases = azurerm_mssql_failover_group.test.databases
  tags      = azurerm_mssql_failover_group.test.tags

  partner_server {
    id = azurerm_mssql_failover_group.test.partner_server[0].id
  }

  read_write_endpoint_failover_policy {
    mode          = azurerm_mssql_failover_group.test.read_write_endpoint_failover_policy[0].mode
    grace_minutes = azurerm_mssql_failover_group.test.read_write_endpoint_failover_policy[0].grace_minutes
  }
}
`, r.automaticFailoverWithDatabases(data))
}
