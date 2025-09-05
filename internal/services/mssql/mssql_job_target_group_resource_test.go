// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobtargetgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MsSqlJobTargetGroupResource struct{}

func TestAccMsSqlJobTargetGroupTest_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_target_group", "test")
	r := MsSqlJobTargetGroupResource{}

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

func TestAccMsSqlJobTargetGroupTest_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_target_group", "test")
	r := MsSqlJobTargetGroupResource{}

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

func TestAccMsSqlJobTargetGroupTest_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_target_group", "test")
	r := MsSqlJobTargetGroupResource{}

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

func TestAccMsSqlJobTargetGroupTest_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_target_group", "test")
	r := MsSqlJobTargetGroupResource{}

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

func TestAccMsSqlJobTargetGroupTest_withDatabase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_target_group", "test")
	r := MsSqlJobTargetGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDatabase(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlJobTargetGroupTest_withElasticPool(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_job_target_group", "test")
	r := MsSqlJobTargetGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withElasticPool(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (MsSqlJobTargetGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := jobtargetgroups.ParseTargetGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.JobTargetGroupsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r MsSqlJobTargetGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_job_target_group" "test" {
  name         = "acctest-target-group-%[2]d"
  job_agent_id = azurerm_mssql_job_agent.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlJobTargetGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_job_target_group" "import" {
  name         = azurerm_mssql_job_target_group.test.name
  job_agent_id = azurerm_mssql_job_target_group.test.job_agent_id
}
`, r.basic(data))
}

func (r MsSqlJobTargetGroupResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_job_target_group" "test" {
  name         = "acctest-target-group-%[2]d"
  job_agent_id = azurerm_mssql_job_agent.test.id

  job_target {
    membership_type   = "Include"
    server_name       = azurerm_mssql_server.test.name
    job_credential_id = azurerm_mssql_job_credential.test.id
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlJobTargetGroupResource) withDatabase(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test2" {
  name      = "acctest-db2-%[2]d"
  server_id = azurerm_mssql_server.test.id
  collation = "SQL_Latin1_General_CP1_CI_AS"
  sku_name  = "S1"
}

resource "azurerm_mssql_job_target_group" "test" {
  name         = "acctest-target-group-%[2]d"
  job_agent_id = azurerm_mssql_job_agent.test.id

  job_target {
    server_name   = azurerm_mssql_server.test.name
    database_name = azurerm_mssql_database.test2.name
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlJobTargetGroupResource) withElasticPool(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_mssql_server.test.name
  max_size_gb         = 5

  sku {
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
    capacity = 4
    family   = "Gen5"
  }

  per_database_settings {
    min_capacity = 0.25
    max_capacity = 4
  }
}

resource "azurerm_mssql_job_target_group" "test" {
  name         = "acctest-target-group-%[2]d"
  job_agent_id = azurerm_mssql_job_agent.test.id

  job_target {
    server_name       = azurerm_mssql_server.test.name
    elastic_pool_name = azurerm_mssql_elasticpool.test.name
    job_credential_id = azurerm_mssql_job_credential.test.id
  }
}
`, r.template(data), data.RandomInteger)
}

func (MsSqlJobTargetGroupResource) template(data acceptance.TestData) string {
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

resource "azurerm_mssql_job_credential" "test" {
  name         = "acctest-job-credential-%[1]d"
  job_agent_id = azurerm_mssql_job_agent.test.id
  username     = "testusername"
  password     = "testpassword"
}
`, data.RandomInteger, data.Locations.Primary)
}
