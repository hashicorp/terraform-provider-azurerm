package streamanalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type StreamAnalyticsStreamInputIoTHubResource struct{}

func TestAccStreamAnalyticsStreamInputIoTHub_avro(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_iothub", "test")
	r := StreamAnalyticsStreamInputIoTHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.avro(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_policy_key"),
	})
}

func TestAccStreamAnalyticsStreamInputIoTHub_csv(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_iothub", "test")
	r := StreamAnalyticsStreamInputIoTHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.csv(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_policy_key"),
	})
}

func TestAccStreamAnalyticsStreamInputIoTHub_json(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_iothub", "test")
	r := StreamAnalyticsStreamInputIoTHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.json(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_policy_key"),
	})
}

func TestAccStreamAnalyticsStreamInputIoTHub_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_iothub", "test")
	r := StreamAnalyticsStreamInputIoTHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.json(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_policy_key"),
	})
}

func TestAccStreamAnalyticsStreamInputIoTHub_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_iothub", "test")
	r := StreamAnalyticsStreamInputIoTHubResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.json(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r StreamAnalyticsStreamInputIoTHubResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	name := state.Attributes["name"]
	jobName := state.Attributes["stream_analytics_job_name"]
	resourceGroup := state.Attributes["resource_group_name"]

	resp, err := client.StreamAnalytics.InputsClient.Get(ctx, resourceGroup, jobName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Stream Output %q (Stream Analytics Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r StreamAnalyticsStreamInputIoTHubResource) avro(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_stream_input_iothub" "test" {
  name                         = "acctestinput-%d"
  stream_analytics_job_name    = azurerm_stream_analytics_job.test.name
  resource_group_name          = azurerm_stream_analytics_job.test.resource_group_name
  endpoint                     = "messages/events"
  iothub_namespace             = azurerm_iothub.test.name
  eventhub_consumer_group_name = "$Default"
  shared_access_policy_key     = azurerm_iothub.test.shared_access_policy[0].primary_key
  shared_access_policy_name    = "iothubowner"

  serialization {
    type = "Avro"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsStreamInputIoTHubResource) csv(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_stream_input_iothub" "test" {
  name                         = "acctestinput-%d"
  stream_analytics_job_name    = azurerm_stream_analytics_job.test.name
  resource_group_name          = azurerm_stream_analytics_job.test.resource_group_name
  endpoint                     = "messages/events"
  iothub_namespace             = azurerm_iothub.test.name
  eventhub_consumer_group_name = "$Default"
  shared_access_policy_key     = azurerm_iothub.test.shared_access_policy[0].primary_key
  shared_access_policy_name    = "iothubowner"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsStreamInputIoTHubResource) json(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_stream_input_iothub" "test" {
  name                         = "acctestinput-%d"
  stream_analytics_job_name    = azurerm_stream_analytics_job.test.name
  resource_group_name          = azurerm_stream_analytics_job.test.resource_group_name
  endpoint                     = "messages/events"
  iothub_namespace             = azurerm_iothub.test.name
  eventhub_consumer_group_name = "$Default"
  shared_access_policy_key     = azurerm_iothub.test.shared_access_policy[0].primary_key
  shared_access_policy_name    = "iothubowner"

  serialization {
    type     = "Json"
    encoding = "UTF8"
  }
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsStreamInputIoTHubResource) updated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub" "updated" {
  name                = "acctestiot2-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }
}

resource "azurerm_stream_analytics_stream_input_iothub" "test" {
  name                         = "acctestinput-%d"
  stream_analytics_job_name    = azurerm_stream_analytics_job.test.name
  resource_group_name          = azurerm_stream_analytics_job.test.resource_group_name
  endpoint                     = "messages/events"
  eventhub_consumer_group_name = "$Default"
  iothub_namespace             = azurerm_iothub.updated.name
  shared_access_policy_key     = azurerm_iothub.updated.shared_access_policy[0].primary_key
  shared_access_policy_name    = "iothubowner"

  serialization {
    type = "Avro"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r StreamAnalyticsStreamInputIoTHubResource) requiresImport(data acceptance.TestData) string {
	template := r.json(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_stream_input_iothub" "import" {
  name                         = azurerm_stream_analytics_stream_input_iothub.test.name
  stream_analytics_job_name    = azurerm_stream_analytics_stream_input_iothub.test.stream_analytics_job_name
  resource_group_name          = azurerm_stream_analytics_stream_input_iothub.test.resource_group_name
  endpoint                     = azurerm_stream_analytics_stream_input_iothub.test.endpoint
  eventhub_consumer_group_name = azurerm_stream_analytics_stream_input_iothub.test.eventhub_consumer_group_name
  iothub_namespace             = azurerm_stream_analytics_stream_input_iothub.test.iothub_namespace
  shared_access_policy_key     = azurerm_stream_analytics_stream_input_iothub.test.shared_access_policy_key
  shared_access_policy_name    = azurerm_stream_analytics_stream_input_iothub.test.resource_group_name

  serialization {
    type     = azurerm_stream_analytics_stream_input_iothub.test.serialization.0.type
    encoding = azurerm_stream_analytics_stream_input_iothub.test.serialization.0.encoding
  }
}
`, template)
}

func (r StreamAnalyticsStreamInputIoTHubResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestiothub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
