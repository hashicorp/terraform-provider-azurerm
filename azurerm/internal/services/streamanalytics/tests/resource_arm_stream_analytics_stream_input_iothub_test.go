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

func TestAccAzureRMStreamAnalyticsStreamInputIoTHub_avro(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_iothub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsStreamInputIoTHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputIoTHub_avro(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputIoTHubExists(data.ResourceName),
				),
			},
			data.ImportStep("shared_access_policy_key"),
		},
	})
}

func TestAccAzureRMStreamAnalyticsStreamInputIoTHub_csv(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_iothub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsStreamInputIoTHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputIoTHub_csv(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputIoTHubExists(data.ResourceName),
				),
			},
			data.ImportStep("shared_access_policy_key"),
		},
	})
}

func TestAccAzureRMStreamAnalyticsStreamInputIoTHub_json(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_iothub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsStreamInputIoTHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputIoTHub_json(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputIoTHubExists(data.ResourceName),
				),
			},
			data.ImportStep("shared_access_policy_key"),
		},
	})
}

func TestAccAzureRMStreamAnalyticsStreamInputIoTHub_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_iothub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsStreamInputIoTHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputIoTHub_json(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputIoTHubExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputIoTHub_updated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputIoTHubExists(data.ResourceName),
				),
			},
			data.ImportStep("shared_access_policy_key"),
		},
	})
}

func TestAccAzureRMStreamAnalyticsStreamInputIoTHub_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_stream_input_iothub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsStreamInputIoTHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputIoTHub_json(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputIoTHubExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMStreamAnalyticsStreamInputIoTHub_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_stream_analytics_stream_input_iothub"),
			},
		},
	})
}

func testCheckAzureRMStreamAnalyticsStreamInputIoTHubExists(resourceName string) resource.TestCheckFunc {
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

func testCheckAzureRMStreamAnalyticsStreamInputIoTHubDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).StreamAnalytics.InputsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_stream_analytics_stream_input_iothub" {
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
			return fmt.Errorf("Stream Analytics Stream Input IoTHub still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMStreamAnalyticsStreamInputIoTHub_avro(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsStreamInputIoTHub_template(data)
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

func testAccAzureRMStreamAnalyticsStreamInputIoTHub_csv(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsStreamInputIoTHub_template(data)
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

func testAccAzureRMStreamAnalyticsStreamInputIoTHub_json(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsStreamInputIoTHub_template(data)
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

func testAccAzureRMStreamAnalyticsStreamInputIoTHub_updated(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsStreamInputIoTHub_template(data)
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

func testAccAzureRMStreamAnalyticsStreamInputIoTHub_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsStreamInputIoTHub_json(data)
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

func testAccAzureRMStreamAnalyticsStreamInputIoTHub_template(data acceptance.TestData) string {
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
