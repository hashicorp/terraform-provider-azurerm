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

type StreamAnalyticsOutputBlobResource struct{}

func TestAccStreamAnalyticsOutputBlob_avro(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_blob", "test")
	r := StreamAnalyticsOutputBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.avro(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
	})
}

func TestAccStreamAnalyticsOutputBlob_csv(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_blob", "test")
	r := StreamAnalyticsOutputBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.csv(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
	})
}

func TestAccStreamAnalyticsOutputBlob_json(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_blob", "test")
	r := StreamAnalyticsOutputBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.json(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
	})
}

func TestAccStreamAnalyticsOutputBlob_parquet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_blob", "test")
	r := StreamAnalyticsOutputBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.parquet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
	})
}

func TestAccStreamAnalyticsOutputBlob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_blob", "test")
	r := StreamAnalyticsOutputBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.json(data),
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

func TestAccStreamAnalyticsOutputBlob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_blob", "test")
	r := StreamAnalyticsOutputBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.json(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccStreamAnalyticsOutputBlob_authenticationMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_blob", "test")
	r := StreamAnalyticsOutputBlobResource{}
	identity := "identity { type = \"SystemAssigned\" }"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.csv(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
		{
			Config: r.authenticationMode(data, identity),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
	})
}

func TestAccStreamAnalyticsOutputBlob_blobWriteMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_blob", "test")
	r := StreamAnalyticsOutputBlobResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blobWriteMode(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
	})
}

func (r StreamAnalyticsOutputBlobResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r StreamAnalyticsOutputBlobResource) avro(data acceptance.TestData) string {
	template := r.template(data, "")
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_blob" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  storage_account_name      = azurerm_storage_account.test.name
  storage_account_key       = azurerm_storage_account.test.primary_access_key
  storage_container_name    = azurerm_storage_container.test.name
  path_pattern              = "some-other-pattern"
  date_format               = "yyyy-MM-dd"
  time_format               = "HH"

  serialization {
    type = "Avro"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputBlobResource) csv(data acceptance.TestData) string {
	template := r.template(data, "")
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_blob" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  storage_account_name      = azurerm_storage_account.test.name
  storage_account_key       = azurerm_storage_account.test.primary_access_key
  storage_container_name    = azurerm_storage_container.test.name
  path_pattern              = "some-pattern"
  date_format               = "yyyy-MM-dd"
  time_format               = "HH"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputBlobResource) json(data acceptance.TestData) string {
	template := r.template(data, "")
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_blob" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  storage_account_name      = azurerm_storage_account.test.name
  storage_account_key       = azurerm_storage_account.test.primary_access_key
  storage_container_name    = azurerm_storage_container.test.name
  path_pattern              = "some-pattern"
  date_format               = "yyyy-MM-dd"
  time_format               = "HH"

  serialization {
    type     = "Json"
    encoding = "UTF8"
    format   = "LineSeparated"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputBlobResource) parquet(data acceptance.TestData) string {
	template := r.template(data, "")
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_blob" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  storage_account_name      = azurerm_storage_account.test.name
  storage_account_key       = azurerm_storage_account.test.primary_access_key
  storage_container_name    = azurerm_storage_container.test.name
  path_pattern              = "some-other-pattern"
  date_format               = "yyyy-MM-dd"
  time_format               = "HH"
  batch_max_wait_time       = "00:02:00"
  batch_min_rows            = 1000000

  serialization {
    type = "Parquet"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputBlobResource) updated(data acceptance.TestData) string {
	template := r.template(data, "")
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "updated" {
  name                     = "acctestsa2%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "updated" {
  name                  = "example"
  storage_account_name  = azurerm_storage_account.updated.name
  container_access_type = "private"
}

resource "azurerm_stream_analytics_output_blob" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  storage_account_name      = azurerm_storage_account.updated.name
  storage_account_key       = azurerm_storage_account.updated.primary_access_key
  storage_container_name    = azurerm_storage_container.updated.name
  path_pattern              = "some-other-pattern"
  date_format               = "yyyy-MM-dd"
  time_format               = "HH"

  serialization {
    type = "Avro"
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func (r StreamAnalyticsOutputBlobResource) requiresImport(data acceptance.TestData) string {
	template := r.json(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_blob" "import" {
  name                      = azurerm_stream_analytics_output_blob.test.name
  stream_analytics_job_name = azurerm_stream_analytics_output_blob.test.stream_analytics_job_name
  resource_group_name       = azurerm_stream_analytics_output_blob.test.resource_group_name
  storage_account_name      = azurerm_stream_analytics_output_blob.test.storage_account_name
  storage_account_key       = azurerm_stream_analytics_output_blob.test.storage_account_key
  storage_container_name    = azurerm_stream_analytics_output_blob.test.storage_container_name
  path_pattern              = azurerm_stream_analytics_output_blob.test.path_pattern
  date_format               = azurerm_stream_analytics_output_blob.test.date_format
  time_format               = azurerm_stream_analytics_output_blob.test.time_format
  dynamic "serialization" {
    for_each = azurerm_stream_analytics_output_blob.test.serialization
    content {
      encoding = lookup(serialization.value, "encoding", null)
      format   = lookup(serialization.value, "format", null)
      type     = serialization.value.type
    }
  }
}
`, template)
}

func (r StreamAnalyticsOutputBlobResource) authenticationMode(data acceptance.TestData, identity string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_blob" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  storage_account_name      = azurerm_storage_account.test.name
  storage_container_name    = azurerm_storage_container.test.name
  path_pattern              = "some-pattern"
  date_format               = "yyyy-MM-dd"
  time_format               = "HH"
  authentication_mode       = "Msi"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}
`, r.template(data, identity), data.RandomInteger)
}

func (r StreamAnalyticsOutputBlobResource) template(data acceptance.TestData, identity string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "example"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_stream_analytics_job" "test" {
  name                                     = "acctestjob-%d"
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

  %s
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, identity)
}

func (r StreamAnalyticsOutputBlobResource) blobWriteMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_blob" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  storage_account_name      = azurerm_storage_account.test.name
  storage_account_key       = azurerm_storage_account.test.primary_access_key
  storage_container_name    = azurerm_storage_container.test.name
  path_pattern              = "some-pattern/{date}/{time}"
  date_format               = "yyyy-MM-dd"
  time_format               = "HH"
  blob_write_mode           = "Once"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}
`, r.template(data, ""), data.RandomInteger)
}
