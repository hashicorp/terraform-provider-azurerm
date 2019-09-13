package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMStreamAnalyticsOutputServiceBusQueue_avro(t *testing.T) {
	resourceName := "azurerm_stream_analytics_output_servicebus_queue.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputServiceBusQueue_avro(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputServiceBusQueueExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"shared_access_policy_key",
				},
			},
		},
	})
}

func TestAccAzureRMStreamAnalyticsOutputServiceBusQueue_csv(t *testing.T) {
	resourceName := "azurerm_stream_analytics_output_servicebus_queue.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputServiceBusQueue_csv(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputServiceBusQueueExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"shared_access_policy_key",
				},
			},
		},
	})
}

func TestAccAzureRMStreamAnalyticsOutputServiceBusQueue_json(t *testing.T) {
	resourceName := "azurerm_stream_analytics_output_servicebus_queue.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputServiceBusQueue_json(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputServiceBusQueueExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"shared_access_policy_key",
				},
			},
		},
	})
}

func TestAccAzureRMStreamAnalyticsOutputServiceBusQueue_update(t *testing.T) {
	resourceName := "azurerm_stream_analytics_output_servicebus_queue.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputServiceBusQueue_json(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputServiceBusQueueExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMStreamAnalyticsOutputServiceBusQueue_updated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputServiceBusQueueExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"shared_access_policy_key",
				},
			},
		},
	})
}

func TestAccAzureRMStreamAnalyticsOutputServiceBusQueue_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_stream_analytics_output_servicebus_queue.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputServiceBusQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputServiceBusQueue_json(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputServiceBusQueueExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMStreamAnalyticsOutputServiceBusQueue_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_stream_analytics_output_servicebus_queue"),
			},
		},
	})
}

func testCheckAzureRMStreamAnalyticsOutputServiceBusQueueExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		jobName := rs.Primary.Attributes["stream_analytics_job_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := testAccProvider.Meta().(*ArmClient).StreamAnalytics.OutputsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, jobName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on streamAnalyticsOutputsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Stream Output %q (Stream Analytics Job %q / Resource Group %q) does not exist", name, jobName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMStreamAnalyticsOutputServiceBusQueueDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).StreamAnalytics.OutputsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_stream_analytics_output_servicebus_queue" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		jobName := rs.Primary.Attributes["stream_analytics_job_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, jobName, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Stream Analytics Output ServiceBus Queue still exists:\n%#v", resp.OutputProperties)
		}
	}

	return nil
}

func testAccAzureRMStreamAnalyticsOutputServiceBusQueue_avro(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputServiceBusQueue_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_servicebus_queue" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name       = "${azurerm_stream_analytics_job.test.resource_group_name}"
  queue_name                = "${azurerm_servicebus_queue.test.name}"
  servicebus_namespace      = "${azurerm_servicebus_namespace.test.name}"
  shared_access_policy_key  = "${azurerm_servicebus_namespace.test.default_primary_key}"
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type = "Avro"
  }
}
`, template, rInt)
}

func testAccAzureRMStreamAnalyticsOutputServiceBusQueue_csv(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputServiceBusQueue_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_servicebus_queue" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name       = "${azurerm_stream_analytics_job.test.resource_group_name}"
  queue_name                = "${azurerm_servicebus_queue.test.name}"
  servicebus_namespace      = "${azurerm_servicebus_namespace.test.name}"
  shared_access_policy_key  = "${azurerm_servicebus_namespace.test.default_primary_key}"
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}
`, template, rInt)
}

func testAccAzureRMStreamAnalyticsOutputServiceBusQueue_json(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputServiceBusQueue_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_servicebus_queue" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name       = "${azurerm_stream_analytics_job.test.resource_group_name}"
  queue_name                = "${azurerm_servicebus_queue.test.name}"
  servicebus_namespace      = "${azurerm_servicebus_namespace.test.name}"
  shared_access_policy_key  = "${azurerm_servicebus_namespace.test.default_primary_key}"
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type     = "Json"
    encoding = "UTF8"
    format   = "LineSeparated"
  }
}
`, template, rInt)
}

func testAccAzureRMStreamAnalyticsOutputServiceBusQueue_updated(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputServiceBusQueue_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_namespace" "updated" {
  name                = "acctest2-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "updated" {
  name                = "acctest2-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.updated.name}"
  enable_partitioning = true
}

resource "azurerm_stream_analytics_output_servicebus_queue" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name       = "${azurerm_stream_analytics_job.test.resource_group_name}"
  queue_name                = "${azurerm_servicebus_queue.updated.name}"
  servicebus_namespace      = "${azurerm_servicebus_namespace.updated.name}"
  shared_access_policy_key  = "${azurerm_servicebus_namespace.updated.default_primary_key}"
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type = "Avro"
  }
}
`, template, rInt, rInt, rInt)
}

func testAccAzureRMStreamAnalyticsOutputServiceBusQueue_requiresImport(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputServiceBusQueue_json(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_servicebus_queue" "import" {
  name                      = "${azurerm_stream_analytics_output_servicebus_queue.test.name}"
  stream_analytics_job_name = "${azurerm_stream_analytics_output_servicebus_queue.test.stream_analytics_job_name}"
  resource_group_name       = "${azurerm_stream_analytics_output_servicebus_queue.test.resource_group_name}"
  queue_name                = "${azurerm_stream_analytics_output_servicebus_queue.test.queue_name}"
  servicebus_namespace      = "${azurerm_stream_analytics_output_servicebus_queue.test.servicebus_namespace}"
  shared_access_policy_key  = "${azurerm_stream_analytics_output_servicebus_queue.test.shared_access_policy_key}"
  shared_access_policy_name = "${azurerm_stream_analytics_output_servicebus_queue.test.shared_access_policy_name}"
  serialization             = "${azurerm_stream_analytics_output_servicebus_queue.test.serialization}"
}
`, template)
}

func testAccAzureRMStreamAnalyticsOutputServiceBusQueue_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  enable_partitioning = true
}

resource "azurerm_stream_analytics_job" "test" {
  name                                     = "acctestjob-%d"
  resource_group_name                      = "${azurerm_resource_group.test.name}"
  location                                 = "${azurerm_resource_group.test.location}"
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
`, rInt, location, rInt, rInt, rInt)
}
