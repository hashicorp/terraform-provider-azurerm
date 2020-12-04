package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMStreamAnalyticsStreamInputEventHub_avro(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_eventhub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsStreamInputEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputEventHub_avro(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputEventHubExists(data.ResourceName),
				),
			},
			data.ImportStep("shared_access_policy_key"),
		},
	})
}

func TestAccAzureRMStreamAnalyticsStreamInputEventHub_csv(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_eventhub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsStreamInputEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputEventHub_csv(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputEventHubExists(data.ResourceName),
				),
			},
			data.ImportStep("shared_access_policy_key"),
		},
	})
}

func TestAccAzureRMStreamAnalyticsStreamInputEventHub_json(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_eventhub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsStreamInputEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputEventHub_json(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputEventHubExists(data.ResourceName),
				),
			},
			data.ImportStep("shared_access_policy_key"),
		},
	})
}

func TestAccAzureRMStreamAnalyticsStreamInputEventHub_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_eventhub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsStreamInputEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputEventHub_json(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputEventHubExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputEventHub_updated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputEventHubExists(data.ResourceName),
				),
			},
			data.ImportStep("shared_access_policy_key"),
		},
	})
}

func TestAccAzureRMStreamAnalyticsStreamInputEventHub_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_eventhub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsStreamInputEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputEventHub_json(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputEventHubExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStreamAnalyticsStreamInputEventHub_requiresImport),
		},
	})
}

func testCheckAzureRMStreamAnalyticsStreamInputEventHubExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).StreamAnalytics.InputsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		jobName := rs.Primary.Attributes["stream_analytics_job_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, jobName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on streamAnalyticsInputsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Stream Input %q (Stream Analytics Job %q / Resource Group %q) does not exist", name, jobName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMStreamAnalyticsStreamInputEventHubDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).StreamAnalytics.InputsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_stream_analytics_stream_input_eventhub" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		jobName := rs.Primary.Attributes["stream_analytics_job_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		resp, err := conn.Get(ctx, resourceGroup, jobName, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Stream Analytics Stream Input EventHub still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMStreamAnalyticsStreamInputEventHub_avro(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsStreamInputEventHub_template(data)
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

  serialization {
    type = "Avro"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMStreamAnalyticsStreamInputEventHub_csv(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsStreamInputEventHub_template(data)
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

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMStreamAnalyticsStreamInputEventHub_json(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsStreamInputEventHub_template(data)
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

  serialization {
    type     = "Json"
    encoding = "UTF8"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMStreamAnalyticsStreamInputEventHub_updated(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsStreamInputEventHub_template(data)
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

  serialization {
    type = "Avro"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMStreamAnalyticsStreamInputEventHub_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsStreamInputEventHub_json(data)
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

func testAccAzureRMStreamAnalyticsStreamInputEventHub_template(data acceptance.TestData) string {
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
