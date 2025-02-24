// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobsteps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MsSqlJobStepTestResource struct{}

func TestAccMsSqlJobStep_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_step", "test")
	r := MsSqlJobStepTestResource{}

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

func TestAccMsSqlJobStep_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_step", "test")
	r := MsSqlJobStepTestResource{}

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

func TestAccMsSqlJobStep_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_step", "test")
	r := MsSqlJobStepTestResource{}

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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlJobStep_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_step", "test")
	r := MsSqlJobStepTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (MsSqlJobStepTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := jobsteps.ParseStepID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.JobStepsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r MsSqlJobStepTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_job_step" "test" {
  name                = "acctest-job-step-%[2]d"
  job_id              = azurerm_mssql_job.test.id
  job_credential_id   = azurerm_mssql_job_credential.test.id
  job_target_group_id = azurerm_mssql_job_target_group.test.id

  job_step_index = 1
  sql_script     = <<EOT
IF NOT EXISTS (SELECT * FROM sys.objects WHERE [name] = N'Person')
  CREATE TABLE Person (
    FirstName NVARCHAR(50),
    LastName NVARCHAR(50),
  );
EOT
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlJobStepTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_job_step" "import" {
  name                = azurerm_mssql_job_step.test.name
  job_id              = azurerm_mssql_job.test.id
  job_credential_id   = azurerm_mssql_job_credential.test.id
  job_target_group_id = azurerm_mssql_job_target_group.test.id

  job_step_index = 1
  sql_script     = <<EOT
IF NOT EXISTS (SELECT * FROM sys.objects WHERE [name] = N'Person')
  CREATE TABLE Person (
    FirstName NVARCHAR(50),
    LastName NVARCHAR(50),
  );
EOT
}
`, r.basic(data))
}

func (r MsSqlJobStepTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_job_step" "test" {
  name                = "acctest-job-step-%[2]d"
  job_id              = azurerm_mssql_job.test.id
  job_credential_id   = azurerm_mssql_job_credential.test.id
  job_target_group_id = azurerm_mssql_job_target_group.test.id

  job_step_index = 1
  sql_script     = <<EOT
IF NOT EXISTS (SELECT * FROM sys.objects WHERE [name] = N'Pets')
  CREATE TABLE Pets (
    Animal NVARCHAR(50),
    Name NVARCHAR(50),
  );
EOT

  initial_retry_interval_seconds    = 1
  maximum_retry_interval_seconds    = 2
  retry_attempts                    = 3
  retry_interval_backoff_multiplier = 4.5
  timeout_seconds                   = 12345

  output_target {
    job_credential_id = azurerm_mssql_job_credential.test.id
    mssql_database_id = azurerm_mssql_database.test.id
    table_name        = "test"
    schema_name       = "dbo"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlJobStepTestResource) template(data acceptance.TestData) string {
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

resource "azurerm_mssql_job" "test" {
  name         = "acctest-job-%[1]d"
  job_agent_id = azurerm_mssql_job_agent.test.id
}

resource "azurerm_mssql_job_target_group" "test" {
  name         = "acctest-target-group-%[1]d"
  job_agent_id = azurerm_mssql_job_agent.test.id
}

resource "azurerm_mssql_job_credential" "test" {
  name         = "acctest-job-credential-%[1]d"
  job_agent_id = azurerm_mssql_job_agent.test.id
  username     = "testusername"
  password     = "testpassword"
}
`, data.RandomInteger, data.Locations.Primary)
}
