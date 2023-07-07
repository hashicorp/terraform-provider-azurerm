// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StreamAnalyticsOutputTableResource struct{}

func TestAccStreamAnalyticsOutputTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_table", "test")
	r := StreamAnalyticsOutputTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
	})
}

func TestAccStreamAnalyticsOutputTable_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_table", "test")
	r := StreamAnalyticsOutputTableResource{}

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
		data.ImportStep("storage_account_key"),
	})
}

func TestAccStreamAnalyticsOutputTable_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_table", "test")
	r := StreamAnalyticsOutputTableResource{}

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

func TestAccStreamAnalyticsOutputTable_columnsToRemove(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_table", "test")
	r := StreamAnalyticsOutputTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
		{
			Config: r.columnsToRemove(data, "column1", "column2", "column3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
		{
			Config: r.columnsToRemove(data, "column2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
	})
}

func (r StreamAnalyticsOutputTableResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r StreamAnalyticsOutputTableResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_table" "test" {
  name                      = "acctestoutput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  storage_account_name      = azurerm_storage_account.test.name
  storage_account_key       = azurerm_storage_account.test.primary_access_key
  table                     = "foobar"
  partition_key             = "foo"
  row_key                   = "bar"
  batch_size                = 100
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputTableResource) updated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "updated" {
  name                     = "acctestaccu%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_table" "updated" {
  name                 = "accteststu%[3]d"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_stream_analytics_output_table" "test" {
  name                      = "acctestoutput-%[3]d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  storage_account_name      = azurerm_storage_account.updated.name
  storage_account_key       = azurerm_storage_account.updated.primary_access_key
  table                     = "updated"
  partition_key             = "partitionkeyupdated"
  row_key                   = "rowkeyupdated"
  batch_size                = 50
}
`, template, data.RandomString, data.RandomInteger)
}

func (r StreamAnalyticsOutputTableResource) columnsToRemove(data acceptance.TestData, columns ...string) string {
	columnsToRemove := make([]string, len(columns))
	for idx, columnToRemove := range columns {
		columnsToRemove[idx] = fmt.Sprintf(`"%s"`, columnToRemove)
	}

	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_table" "test" {
  name                      = "acctestoutput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  storage_account_name      = azurerm_storage_account.test.name
  storage_account_key       = azurerm_storage_account.test.primary_access_key
  table                     = "foobar"
  partition_key             = "foo"
  row_key                   = "bar"
  batch_size                = 100
  columns_to_remove         = [%s]
}
`, template, data.RandomInteger, strings.Join(columnsToRemove, ","))
}

func (r StreamAnalyticsOutputTableResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_table" "import" {
  name                      = azurerm_stream_analytics_output_table.test.name
  stream_analytics_job_name = azurerm_stream_analytics_output_table.test.stream_analytics_job_name
  resource_group_name       = azurerm_stream_analytics_output_table.test.resource_group_name
  storage_account_name      = azurerm_stream_analytics_output_table.test.storage_account_name
  storage_account_key       = azurerm_stream_analytics_output_table.test.storage_account_key
  table                     = azurerm_stream_analytics_output_table.test.table
  partition_key             = azurerm_stream_analytics_output_table.test.partition_key
  row_key                   = azurerm_stream_analytics_output_table.test.row_key
  batch_size                = azurerm_stream_analytics_output_table.test.batch_size
}
`, template)
}

func (r StreamAnalyticsOutputTableResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%[1]d"
  storage_account_name = azurerm_storage_account.test.name
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
