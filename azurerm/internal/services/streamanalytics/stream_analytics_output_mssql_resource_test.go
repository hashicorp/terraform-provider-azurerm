package streamanalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

func (r StreamAnalyticsOutputSqlResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	name := state.Attributes["name"]
	jobName := state.Attributes["stream_analytics_job_name"]
	resourceGroup := state.Attributes["resource_group_name"]

	resp, err := client.StreamAnalytics.OutputsClient.Get(ctx, resourceGroup, jobName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Stream Output %q (Stream Analytics Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
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

  server   = azurerm_sql_server.test.fully_qualified_domain_name
  user     = azurerm_sql_server.test.administrator_login
  password = azurerm_sql_server.test.administrator_login_password
  database = azurerm_sql_database.test.name
  table    = "AccTestTable"
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputSqlResource) updated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_mssql" "test" {
  name                      = "acctestoutput-updated-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name

  server   = azurerm_sql_server.test.fully_qualified_domain_name
  user     = azurerm_sql_server.test.administrator_login
  password = azurerm_sql_server.test.administrator_login_password
  database = azurerm_sql_database.test.name
  table    = "AccTestTable"
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

  server   = azurerm_sql_server.test.fully_qualified_domain_name
  user     = azurerm_sql_server.test.administrator_login
  password = azurerm_sql_server.test.administrator_login_password
  database = azurerm_sql_database.test.name
  table    = "AccTestTable"
}
`, template)
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

resource "azurerm_sql_server" "test" {
  name                         = "acctestserver-%s"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "acctestadmin"
  administrator_login_password = "t2RX8A76GrnE4EKC"
}

resource "azurerm_sql_database" "test" {
  name                             = "acctestdb"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = azurerm_resource_group.test.location
  server_name                      = azurerm_sql_server.test.name
  requested_service_objective_name = "S0"
  collation                        = "SQL_LATIN1_GENERAL_CP1_CI_AS"
  max_size_bytes                   = "268435456000"
  create_mode                      = "Default"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
