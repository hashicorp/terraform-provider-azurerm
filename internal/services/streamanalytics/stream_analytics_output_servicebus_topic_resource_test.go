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

type StreamAnalyticsOutputServiceBusTopicResource struct{}

func TestAccStreamAnalyticsOutputServiceBusTopic_avro(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_servicebus_topic", "test")
	r := StreamAnalyticsOutputServiceBusTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.avro(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_policy_key"),
	})
}

func TestAccStreamAnalyticsOutputServiceBusTopic_csv(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_servicebus_topic", "test")
	r := StreamAnalyticsOutputServiceBusTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.csv(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_policy_key"),
	})
}

func TestAccStreamAnalyticsOutputServiceBusTopic_json(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_servicebus_topic", "test")
	r := StreamAnalyticsOutputServiceBusTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.json(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_policy_key"),
	})
}

func TestAccStreamAnalyticsOutputServiceBusTopic_propertyColumns(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_servicebus_topic", "test")
	r := StreamAnalyticsOutputServiceBusTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.propertyColumns(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("property_columns.0").HasValue("col1"),
				check.That(data.ResourceName).Key("property_columns.1").HasValue("col2"),
			),
		},
		data.ImportStep("shared_access_policy_key"),
	})
}

func TestAccStreamAnalyticsOutputServiceBusTopic_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_servicebus_topic", "test")
	r := StreamAnalyticsOutputServiceBusTopicResource{}

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
		data.ImportStep("shared_access_policy_key"),
	})
}

func TestAccStreamAnalyticsOutputServiceBusTopic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_servicebus_topic", "test")
	r := StreamAnalyticsOutputServiceBusTopicResource{}

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

func TestAccStreamAnalyticsOutputServiceBusTopic_authenticationModeMsi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_servicebus_topic", "test")
	r := StreamAnalyticsOutputServiceBusTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authenticationModeMsi(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_policy_key"),
	})
}

func TestAccStreamAnalyticsOutputServiceBusTopic_systemPropertyColumns(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_output_servicebus_topic", "test")
	r := StreamAnalyticsOutputServiceBusTopicResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.csv(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_policy_key"),
		{
			Config: r.systemPropertyColumns(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_policy_key"),
		{
			Config: r.updateSystemPropertyColumns(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_policy_key"),
		{
			Config: r.csv(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_policy_key"),
	})
}

func (r StreamAnalyticsOutputServiceBusTopicResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := outputs.ParseOutputID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.StreamAnalytics.OutputsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving (%s): %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r StreamAnalyticsOutputServiceBusTopicResource) avro(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_servicebus_topic" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  topic_name                = azurerm_servicebus_topic.test.name
  servicebus_namespace      = azurerm_servicebus_namespace.test.name
  shared_access_policy_key  = azurerm_servicebus_namespace.test.default_primary_key
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type = "Avro"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputServiceBusTopicResource) csv(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_servicebus_topic" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  topic_name                = azurerm_servicebus_topic.test.name
  servicebus_namespace      = azurerm_servicebus_namespace.test.name
  shared_access_policy_key  = azurerm_servicebus_namespace.test.default_primary_key
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputServiceBusTopicResource) json(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_servicebus_topic" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  topic_name                = azurerm_servicebus_topic.test.name
  servicebus_namespace      = azurerm_servicebus_namespace.test.name
  shared_access_policy_key  = azurerm_servicebus_namespace.test.default_primary_key
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type     = "Json"
    encoding = "UTF8"
    format   = "LineSeparated"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputServiceBusTopicResource) authenticationModeMsi(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_servicebus_topic" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  topic_name                = azurerm_servicebus_topic.test.name
  servicebus_namespace      = azurerm_servicebus_namespace.test.name
  authentication_mode       = "Msi"

  serialization {
    type     = "Json"
    encoding = "UTF8"
    format   = "LineSeparated"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputServiceBusTopicResource) propertyColumns(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_servicebus_topic" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  topic_name                = azurerm_servicebus_topic.test.name
  servicebus_namespace      = azurerm_servicebus_namespace.test.name
  shared_access_policy_key  = azurerm_servicebus_namespace.test.default_primary_key
  shared_access_policy_name = "RootManageSharedAccessKey"
  property_columns          = ["col1", "col2"]

  serialization {
    type     = "Json"
    encoding = "UTF8"
    format   = "LineSeparated"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputServiceBusTopicResource) updated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace" "updated" {
  name                = "acctest2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "updated" {
  name                = "acctest2-%d"
  namespace_id        = azurerm_servicebus_namespace.updated.id
  enable_partitioning = true
}

resource "azurerm_stream_analytics_output_servicebus_topic" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  topic_name                = azurerm_servicebus_topic.updated.name
  servicebus_namespace      = azurerm_servicebus_namespace.updated.name
  shared_access_policy_key  = azurerm_servicebus_namespace.updated.default_primary_key
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type = "Avro"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r StreamAnalyticsOutputServiceBusTopicResource) systemPropertyColumns(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_servicebus_topic" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  topic_name                = azurerm_servicebus_topic.test.name
  servicebus_namespace      = azurerm_servicebus_namespace.test.name
  shared_access_policy_key  = azurerm_servicebus_namespace.test.default_primary_key
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }

  system_property_columns = {
    TimeToLive    = "100"
    Label         = "Test"
    CorrelationId = "79b839ac-be78-4542-8185-098170483986"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputServiceBusTopicResource) updateSystemPropertyColumns(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_servicebus_topic" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  topic_name                = azurerm_servicebus_topic.test.name
  servicebus_namespace      = azurerm_servicebus_namespace.test.name
  shared_access_policy_key  = azurerm_servicebus_namespace.test.default_primary_key
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }

  system_property_columns = {
    TimeToLive    = "101"
    CorrelationId = "79b839ac-be78-4542-8185-098170483986"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsOutputServiceBusTopicResource) requiresImport(data acceptance.TestData) string {
	template := r.json(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_servicebus_topic" "import" {
  name                      = azurerm_stream_analytics_output_servicebus_topic.test.name
  stream_analytics_job_name = azurerm_stream_analytics_output_servicebus_topic.test.stream_analytics_job_name
  resource_group_name       = azurerm_stream_analytics_output_servicebus_topic.test.resource_group_name
  topic_name                = azurerm_stream_analytics_output_servicebus_topic.test.topic_name
  servicebus_namespace      = azurerm_stream_analytics_output_servicebus_topic.test.servicebus_namespace
  shared_access_policy_key  = azurerm_stream_analytics_output_servicebus_topic.test.shared_access_policy_key
  shared_access_policy_name = azurerm_stream_analytics_output_servicebus_topic.test.shared_access_policy_name
  dynamic "serialization" {
    for_each = azurerm_stream_analytics_output_servicebus_topic.test.serialization
    content {
      encoding = lookup(serialization.value, "encoding", null)
      format   = lookup(serialization.value, "format", null)
      type     = serialization.value.type
    }
  }
}
`, template)
}

func (r StreamAnalyticsOutputServiceBusTopicResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name                = "acctest-%d"
  namespace_id        = azurerm_servicebus_namespace.test.id
  enable_partitioning = true
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

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
