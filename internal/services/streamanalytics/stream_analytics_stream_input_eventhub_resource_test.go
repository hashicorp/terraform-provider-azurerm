// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StreamAnalyticsStreamInputEventHubResource struct{}

func TestAccStreamAnalyticsStreamInputEventHub_avro(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_eventhub", "test")
	r := StreamAnalyticsStreamInputEventHubResource{}

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

func TestAccStreamAnalyticsStreamInputEventHub_csv(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_eventhub", "test")
	r := StreamAnalyticsStreamInputEventHubResource{}

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

func TestAccStreamAnalyticsStreamInputEventHub_json(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_eventhub", "test")
	r := StreamAnalyticsStreamInputEventHubResource{}

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

func TestAccStreamAnalyticsStreamInputEventHub_noOptional(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_eventhub", "test")
	r := StreamAnalyticsStreamInputEventHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.jsonNoOptional(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("eventhub_consumer_group_name").IsEmpty(),
				check.That((data.ResourceName)).Key("partition_key").IsEmpty(),
			),
		},
		data.ImportStep("shared_access_policy_key"),
	})
}

func TestAccStreamAnalyticsStreamInputEventHub_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_eventhub", "test")
	r := StreamAnalyticsStreamInputEventHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.json(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("eventhub_consumer_group_name").MatchesOtherKey(
					check.That("azurerm_eventhub_consumer_group.test").Key("name")),
				check.That((data.ResourceName)).Key("partition_key").HasValue("partitionKey"),
			),
		},
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("eventhub_consumer_group_name").MatchesOtherKey(
					check.That("azurerm_eventhub_consumer_group.updated").Key("name")),
				check.That((data.ResourceName)).Key("partition_key").HasValue("updatedPartitionKey"),
			),
		},
		{
			Config: r.jsonNoOptional(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("eventhub_consumer_group_name").IsEmpty(),
				check.That((data.ResourceName)).Key("partition_key").IsEmpty(),
			),
		},
		data.ImportStep("shared_access_policy_key"),
	})
}

func TestAccStreamAnalyticsStreamInputEventHub_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_eventhub", "test")
	r := StreamAnalyticsStreamInputEventHubResource{}

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

func TestAccStreamAnalyticsStreamInputEventHub_authenticationMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_eventhub", "test")
	r := StreamAnalyticsStreamInputEventHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authenticationMode(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_policy_key"),
	})
}

func TestAccStreamAnalyticsStreamInputEventHub_msiWithoutSharedAccessPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_eventhub", "test")
	r := StreamAnalyticsStreamInputEventHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.msiWithoutSharedAccessPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_policy_key"),
	})
}

func (r StreamAnalyticsStreamInputEventHubResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := inputs.ParseInputID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.StreamAnalytics.InputsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r StreamAnalyticsStreamInputEventHubResource) avro(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_stream_input_eventhub" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  eventhub_name             = azurerm_eventhub.test.name
  servicebus_namespace      = azurerm_eventhub_namespace.test.name
  shared_access_policy_key  = azurerm_eventhub_namespace.test.default_primary_key
  shared_access_policy_name = "RootManageSharedAccessKey"
  partition_key             = "partitionKey"

  serialization {
    type = "Avro"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsStreamInputEventHubResource) csv(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_stream_input_eventhub" "test" {
  name                         = "acctestinput-%d"
  stream_analytics_job_name    = azurerm_stream_analytics_job.test.name
  resource_group_name          = azurerm_stream_analytics_job.test.resource_group_name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name
  eventhub_name                = azurerm_eventhub.test.name
  servicebus_namespace         = azurerm_eventhub_namespace.test.name
  shared_access_policy_key     = azurerm_eventhub_namespace.test.default_primary_key
  shared_access_policy_name    = "RootManageSharedAccessKey"
  partition_key                = "partitionKey"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsStreamInputEventHubResource) json(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_stream_input_eventhub" "test" {
  name                         = "acctestinput-%d"
  stream_analytics_job_name    = azurerm_stream_analytics_job.test.name
  resource_group_name          = azurerm_stream_analytics_job.test.resource_group_name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name
  eventhub_name                = azurerm_eventhub.test.name
  servicebus_namespace         = azurerm_eventhub_namespace.test.name
  shared_access_policy_key     = azurerm_eventhub_namespace.test.default_primary_key
  shared_access_policy_name    = "RootManageSharedAccessKey"
  partition_key                = "partitionKey"

  serialization {
    type     = "Json"
    encoding = "UTF8"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsStreamInputEventHubResource) jsonNoOptional(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_stream_input_eventhub" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  eventhub_name             = azurerm_eventhub.test.name
  servicebus_namespace      = azurerm_eventhub_namespace.test.name
  shared_access_policy_key  = azurerm_eventhub_namespace.test.default_primary_key
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type     = "Json"
    encoding = "UTF8"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsStreamInputEventHubResource) updated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_namespace" "updated" {
  name                = "acctestehn2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = 1
}

resource "azurerm_eventhub" "updated" {
  name                = "acctesteh2-%d"
  namespace_name      = azurerm_eventhub_namespace.updated.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_consumer_group" "updated" {
  name                = "acctesteventhubcg2-%d"
  namespace_name      = azurerm_eventhub_namespace.updated.name
  eventhub_name       = azurerm_eventhub.updated.name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_stream_analytics_stream_input_eventhub" "test" {
  name                         = "acctestinput-%d"
  stream_analytics_job_name    = azurerm_stream_analytics_job.test.name
  resource_group_name          = azurerm_stream_analytics_job.test.resource_group_name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.updated.name
  eventhub_name                = azurerm_eventhub.updated.name
  servicebus_namespace         = azurerm_eventhub_namespace.updated.name
  shared_access_policy_key     = azurerm_eventhub_namespace.updated.default_primary_key
  shared_access_policy_name    = "RootManageSharedAccessKey"
  partition_key                = "updatedPartitionKey"

  serialization {
    type = "Avro"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r StreamAnalyticsStreamInputEventHubResource) requiresImport(data acceptance.TestData) string {
	template := r.json(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_stream_input_eventhub" "import" {
  name                         = azurerm_stream_analytics_stream_input_eventhub.test.name
  stream_analytics_job_name    = azurerm_stream_analytics_stream_input_eventhub.test.stream_analytics_job_name
  resource_group_name          = azurerm_stream_analytics_stream_input_eventhub.test.resource_group_name
  eventhub_consumer_group_name = azurerm_stream_analytics_stream_input_eventhub.test.eventhub_consumer_group_name
  eventhub_name                = azurerm_stream_analytics_stream_input_eventhub.test.eventhub_name
  servicebus_namespace         = azurerm_stream_analytics_stream_input_eventhub.test.servicebus_namespace
  shared_access_policy_key     = azurerm_stream_analytics_stream_input_eventhub.test.shared_access_policy_key
  shared_access_policy_name    = azurerm_stream_analytics_stream_input_eventhub.test.shared_access_policy_name
  dynamic "serialization" {
    for_each = azurerm_stream_analytics_stream_input_eventhub.test.serialization
    content {
      encoding = lookup(serialization.value, "encoding", null)
      type     = serialization.value.type
    }
  }
}
`, template)
}

func (r StreamAnalyticsStreamInputEventHubResource) authenticationMode(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_stream_input_eventhub" "test" {
  name                         = "acctestinput-%d"
  stream_analytics_job_name    = azurerm_stream_analytics_job.test.name
  resource_group_name          = azurerm_stream_analytics_job.test.resource_group_name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name
  eventhub_name                = azurerm_eventhub.test.name
  servicebus_namespace         = azurerm_eventhub_namespace.test.name
  shared_access_policy_key     = azurerm_eventhub_namespace.test.default_primary_key
  shared_access_policy_name    = "RootManagedSharedAccessKey"
  partition_key                = "partitionKey"
  authentication_mode          = "ConnectionString"

  serialization {
    type     = "Json"
    encoding = "UTF8"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsStreamInputEventHubResource) msiWithoutSharedAccessPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_stream_input_eventhub" "test" {
  name                         = "acctestinput-%d"
  stream_analytics_job_name    = azurerm_stream_analytics_job.test.name
  resource_group_name          = azurerm_stream_analytics_job.test.resource_group_name
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.test.name
  eventhub_name                = azurerm_eventhub.test.name
  servicebus_namespace         = azurerm_eventhub_namespace.test.name
  authentication_mode          = "Msi"

  serialization {
    type     = "Json"
    encoding = "UTF8"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r StreamAnalyticsStreamInputEventHubResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctestehn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = 1
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteh-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_consumer_group" "test" {
  name                = "acctesteventhubcg-%d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
