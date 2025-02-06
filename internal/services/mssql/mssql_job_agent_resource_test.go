// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobagents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MsSqlJobAgentResource struct{}

func TestAccMsSqlJobAgent_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_agent", "test")
	r := MsSqlJobAgentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlJobAgent_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_agent", "test")
	r := MsSqlJobAgentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMsSqlJobAgent_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_agent", "test")
	r := MsSqlJobAgentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (MsSqlJobAgentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := jobagents.ParseJobAgentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.JobAgentsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("%s does not exist", *id)
		}
		return nil, fmt.Errorf("reading %s: %v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (MsSqlJobAgentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-jobagent-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestmssqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dministr4t0r"
  administrator_login_password = "superSecur3!!!"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctestmssqldb%[1]d"
  server_id = azurerm_mssql_server.test.id
  collation = "SQL_Latin1_General_CP1_CI_AS"
  sku_name  = "S1"
}

resource "azurerm_mssql_job_agent" "test" {
  name        = "acctestmssqljobagent%[1]d"
  location    = azurerm_resource_group.test.location
  database_id = azurerm_mssql_database.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r MsSqlJobAgentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_job_agent" "import" {
  name        = azurerm_mssql_job_agent.test.name
  location    = azurerm_resource_group.test.location
  database_id = azurerm_mssql_database.test.id
}
`, r.basic(data))
}

func (MsSqlJobAgentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-jobagent-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestmssqlserver%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dministr4t0r"
  administrator_login_password = "superSecur3!!!"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctestmssqldb%[1]d"
  server_id = azurerm_mssql_server.test.id
  collation = "SQL_Latin1_General_CP1_CI_AS"
  sku_name  = "S1"
}

resource "azurerm_mssql_job_agent" "test" {
  name        = "acctestmssqljobagent%[1]d"
  location    = azurerm_resource_group.test.location
  database_id = azurerm_mssql_database.test.id

  tags = {
    ENV = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
