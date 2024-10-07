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

type StreamAnalyticsOutputCosmosDBResource struct{}

func TestAccStreamAnalyticsOutputCosmosDB_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_cosmosdb", "test")
	r := StreamAnalyticsOutputCosmosDBResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_account_key"),
	})
}

func TestAccStreamAnalyticsOutputCosmosDB_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_cosmosdb", "test")
	r := StreamAnalyticsOutputCosmosDBResource{}

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
		data.ImportStep("cosmosdb_account_key"),
	})
}

func TestAccStreamAnalyticsOutputCosmosDB_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_cosmosdb", "test")
	r := StreamAnalyticsOutputCosmosDBResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("cosmosdb_account_key"),
	})
}

func TestAccStreamAnalyticsOutputCosmosDB_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_cosmosdb", "test")
	r := StreamAnalyticsOutputCosmosDBResource{}

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
func (r StreamAnalyticsOutputCosmosDBResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r StreamAnalyticsOutputCosmosDBResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_stream_analytics_output_cosmosdb" "test" {
  name                     = "acctestoutput-%[3]d"
  stream_analytics_job_id  = azurerm_stream_analytics_job.test.id
  cosmosdb_account_key     = azurerm_cosmosdb_account.test.primary_key
  cosmosdb_sql_database_id = azurerm_cosmosdb_sql_database.test.id
  container_name           = azurerm_cosmosdb_sql_container.test.name
}
`, template, data.RandomString, data.RandomInteger)
}

func (r StreamAnalyticsOutputCosmosDBResource) updated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_database" "updated" {
  name                = "updated-cosmos-sql-db"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  throughput          = 400
}

resource "azurerm_cosmosdb_sql_container" "updated" {
  name                = "updated-container%[2]s"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.updated.name
  partition_key_paths = ["/definition/id"]
}

resource "azurerm_stream_analytics_output_cosmosdb" "test" {
  name                     = "acctestoutput-%[3]d"
  stream_analytics_job_id  = azurerm_stream_analytics_job.test.id
  cosmosdb_account_key     = azurerm_cosmosdb_account.test.primary_key
  cosmosdb_sql_database_id = azurerm_cosmosdb_sql_database.updated.id
  container_name           = azurerm_cosmosdb_sql_container.updated.name
}
`, template, data.RandomString, data.RandomInteger)
}

func (r StreamAnalyticsOutputCosmosDBResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_stream_analytics_output_cosmosdb" "test" {
  name                     = "acctestoutput-%[3]d"
  stream_analytics_job_id  = azurerm_stream_analytics_job.test.id
  cosmosdb_account_key     = azurerm_cosmosdb_account.test.primary_key
  cosmosdb_sql_database_id = azurerm_cosmosdb_sql_database.test.id
  container_name           = azurerm_cosmosdb_sql_container.test.name
  document_id              = "exampledocumentid"
  partition_key            = "examplekey"
}
`, template, data.RandomString, data.RandomInteger)
}

func (r StreamAnalyticsOutputCosmosDBResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_cosmosdb" "import" {
  name                     = azurerm_stream_analytics_output_cosmosdb.test.name
  stream_analytics_job_id  = azurerm_stream_analytics_output_cosmosdb.test.stream_analytics_job_id
  cosmosdb_account_key     = azurerm_stream_analytics_output_cosmosdb.test.cosmosdb_account_key
  cosmosdb_sql_database_id = azurerm_stream_analytics_output_cosmosdb.test.cosmosdb_sql_database_id
  container_name           = azurerm_stream_analytics_output_cosmosdb.test.container_name
}
`, template)
}

func (r StreamAnalyticsOutputCosmosDBResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctestacc%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 10
    max_staleness_prefix    = 200
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "cosmos-sql-db"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  throughput          = 400
}

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "test-container%[2]s"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition/id"]
}

resource "azurerm_stream_analytics_job" "test" {
  name                                     = "acctestjob-%[1]d"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
