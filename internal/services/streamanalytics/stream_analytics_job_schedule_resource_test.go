// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StreamAnalyticsJobScheduleResource struct{}

func TestAccStreamAnalyticsJobSchedule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job_schedule", "test")
	r := StreamAnalyticsJobScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// todo framework
		// `last_output_time` has different values between refresh steps so we'll ignore it until framework goes in
		data.ImportStep("last_output_time"),
	})
}

func TestAccStreamAnalyticsJobSchedule_customTime(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job_schedule", "test")
	r := StreamAnalyticsJobScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customTime(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// todo framework
		// `last_output_time` has different values between refresh steps so we'll ignore it until framework goes in
		data.ImportStep("last_output_time"),
	})
}

func TestAccStreamAnalyticsJobSchedule_lastOutputEventTime(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job_schedule", "test")
	r := StreamAnalyticsJobScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// todo framework
		// `last_output_time` has different values between refresh steps so we'll ignore it until framework goes in
		data.ImportStep("last_output_time"),
		{
			Config: r.lastOutputEventTime(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// todo framework
		// `last_output_time` has different values between refresh steps so we'll ignore it until framework goes in
		data.ImportStep("last_output_time"),
	})
}

func (r StreamAnalyticsJobScheduleResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.StreamingJobScheduleID(state.ID)
	if err != nil {
		return nil, err
	}

	streamingJobId := streamingjobs.NewStreamingJobID(id.SubscriptionId, id.ResourceGroup, id.StreamingJobName)

	var opts streamingjobs.GetOperationOptions
	resp, err := client.StreamAnalytics.JobsClient.Get(ctx, streamingJobId, opts)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, err
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(resp.Model != nil && resp.Model.Properties.OutputStartTime != nil), nil
}

func (r StreamAnalyticsJobScheduleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_job_schedule" "test" {
  stream_analytics_job_id = azurerm_stream_analytics_job.test.id
  start_mode              = "JobStartTime"

  depends_on = [
    azurerm_stream_analytics_job.test,
    azurerm_stream_analytics_stream_input_blob.test,
    azurerm_stream_analytics_output_blob.test,
  ]
}
`, r.template(data))
}

func (r StreamAnalyticsJobScheduleResource) customTime(data acceptance.TestData) string {
	utcNow := time.Now().UTC()
	startDate := time.Date(utcNow.Year(), utcNow.Month(), 1, 0, 0, 0, 0, utcNow.Location())

	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_job_schedule" "test" {
  stream_analytics_job_id = azurerm_stream_analytics_job.test.id
  start_mode              = "CustomTime"
  start_time              = "%s"

  depends_on = [
    azurerm_stream_analytics_job.test,
    azurerm_stream_analytics_stream_input_blob.test,
    azurerm_stream_analytics_output_blob.test,
  ]
}
`, r.template(data), startDate.Format(time.RFC3339))
}

func (r StreamAnalyticsJobScheduleResource) lastOutputEventTime(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_job_schedule" "test" {
  stream_analytics_job_id = azurerm_stream_analytics_job.test.id
  start_mode              = "LastOutputEventTime"

  depends_on = [
    azurerm_stream_analytics_job.test,
    azurerm_stream_analytics_stream_input_blob.test,
    azurerm_stream_analytics_output_blob.test,
  ]
}


`, r.template(data))
}

func (r StreamAnalyticsJobScheduleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "chonks"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test" {
  name                   = "chonkdata"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source                 = "testdata/chonkdata.csv"
}

resource "azurerm_stream_analytics_job" "test" {
  name                                     = "acctestjob-%[1]d"
  resource_group_name                      = azurerm_resource_group.test.name
  location                                 = azurerm_resource_group.test.location
  data_locale                              = "en-GB"
  compatibility_level                      = "1.1"
  events_late_arrival_max_delay_in_seconds = 10
  events_out_of_order_max_delay_in_seconds = 20
  events_out_of_order_policy               = "Drop"
  output_error_policy                      = "Stop"
  streaming_units                          = 6

  transformation_query = <<QUERY
    SELECT *
    INTO [acctestoutputchonk]
    FROM [acctestinputchonk]
QUERY
}

resource "azurerm_stream_analytics_stream_input_blob" "test" {
  name                      = "acctestinputchonk"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  storage_account_name      = azurerm_storage_account.test.name
  storage_account_key       = azurerm_storage_account.test.primary_access_key
  storage_container_name    = azurerm_storage_container.test.name
  path_pattern              = ""
  date_format               = "yyyy/MM/dd"
  time_format               = "HH"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}

resource "azurerm_stream_analytics_output_blob" "test" {
  name                      = "acctestoutputchonk"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  storage_account_name      = azurerm_storage_account.test.name
  storage_account_key       = azurerm_storage_account.test.primary_access_key
  storage_container_name    = azurerm_storage_container.test.name
  path_pattern              = "avro-chonks-{date}-{time}"
  date_format               = "yyyy-MM-dd"
  time_format               = "HH"

  serialization {
    type = "Avro"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
