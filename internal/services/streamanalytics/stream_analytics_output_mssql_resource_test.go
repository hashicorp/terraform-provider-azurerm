// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StreamAnalyticsOutputSqlResource struct{}

func TestAccStreamAnalyticsOutputSql_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_mssql", "test")
	r := StreamAnalyticsOutputSqlResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccStreamAnalyticsOutputSql_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_mssql", "test")
	r := StreamAnalyticsOutputSqlResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccStreamAnalyticsOutputSql_authenticationModeMsi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_mssql", "test")
	r := StreamAnalyticsOutputSqlResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authenticationModeMsi(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccStreamAnalyticsOutputSql_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_mssql", "test")
	r := StreamAnalyticsOutputSqlResource{}

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

func TestAccStreamAnalyticsOutputSql_maxBatchCountAndMaxWriterCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_mssql", "test")
	r := StreamAnalyticsOutputSqlResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
		{
			Config: r.maxBatchCountAndMaxWriterCount(data, 10001, 0),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
		{
			Config: r.maxBatchCountAndMaxWriterCount(data, 10002, 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func (r StreamAnalyticsOutputSqlResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := outputs.ParseOutputID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.StreamAnalytics.OutputsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r StreamAnalyticsOutputSqlResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_mssql" "test" {
  name                      = "acctestoutput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name

  server   = azurerm_mssql_server.test.fully_qualified_domain_name
  user     = azurerm_mssql_server.test.administrator_login
  password = azurerm_mssql_server.test.administrator_login_password
  database = azurerm_mssql_database.test.name
  table    = "AccTestTable"
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputSqlResource) updated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_mssql" "test" {
  name                      = "acctestoutput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name

  server   = azurerm_mssql_server.test.fully_qualified_domain_name
  user     = azurerm_mssql_server.test.administrator_login
  password = azurerm_mssql_server.test.administrator_login_password
  database = azurerm_mssql_database.test.name
  table    = "AccTestTable"

  max_batch_count = 1000
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputSqlResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_mssql" "import" {
  name                      = azurerm_stream_analytics_output_mssql.test.name
  stream_analytics_job_name = azurerm_stream_analytics_output_mssql.test.stream_analytics_job_name
  resource_group_name       = azurerm_stream_analytics_output_mssql.test.resource_group_name

  server   = azurerm_mssql_server.test.fully_qualified_domain_name
  user     = azurerm_mssql_server.test.administrator_login
  password = azurerm_mssql_server.test.administrator_login_password
  database = azurerm_mssql_database.test.name
  table    = "AccTestTable"
}
`, template)
}

func (r StreamAnalyticsOutputSqlResource) maxBatchCountAndMaxWriterCount(data acceptance.TestData, maxBatchCount, maxWriterCount float64) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_mssql" "test" {
  name                      = "acctestoutput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name

  server   = azurerm_mssql_server.test.fully_qualified_domain_name
  user     = azurerm_mssql_server.test.administrator_login
  password = azurerm_mssql_server.test.administrator_login_password
  database = azurerm_mssql_database.test.name
  table    = "AccTestTable"

  max_batch_count  = %f
  max_writer_count = %f
}
`, template, data.RandomInteger, maxBatchCount, maxWriterCount)
}

func (r StreamAnalyticsOutputSqlResource) authenticationModeMsi(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_mssql" "test" {
  name                      = "acctestoutput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  authentication_mode       = "Msi"

  server   = azurerm_mssql_server.test.fully_qualified_domain_name
  database = azurerm_mssql_database.test.name
  table    = "AccTestTable"
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputSqlResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctestsqlserver%s"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%s"
  server_id = azurerm_mssql_server.test.id
}

resource "azurerm_stream_analytics_job" "test" {
  name                                     = "acctestjob-%s"
  resource_group_name                      = azurerm_resource_group.test.name
  location                                 = azurerm_resource_group.test.location
  compatibility_level                      = "1.0"
  data_locale                              = "en-GB"
  events_late_arrival_max_delay_in_seconds = 60
  events_out_of_order_max_delay_in_seconds = 50
  events_out_of_order_policy               = "Adjust"
  output_error_policy                      = "Drop"
  streaming_units                          = 3

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY

}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
}
