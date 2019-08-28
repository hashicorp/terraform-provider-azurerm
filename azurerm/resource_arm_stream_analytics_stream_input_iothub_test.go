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

func TestAccAzureRMStreamAnalyticsStreamInputIoTHub_avro(t *testing.T) {
	resourceName := "azurerm_stream_analytics_stream_input_iothub.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsStreamInputIoTHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputIoTHub_avro(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputIoTHubExists(resourceName),
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

func TestAccAzureRMStreamAnalyticsStreamInputIoTHub_csv(t *testing.T) {
	resourceName := "azurerm_stream_analytics_stream_input_iothub.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsStreamInputIoTHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputIoTHub_csv(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputIoTHubExists(resourceName),
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

func TestAccAzureRMStreamAnalyticsStreamInputIoTHub_json(t *testing.T) {
	resourceName := "azurerm_stream_analytics_stream_input_iothub.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsStreamInputIoTHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputIoTHub_json(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputIoTHubExists(resourceName),
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

func TestAccAzureRMStreamAnalyticsStreamInputIoTHub_update(t *testing.T) {
	resourceName := "azurerm_stream_analytics_stream_input_iothub.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsStreamInputIoTHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputIoTHub_json(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputIoTHubExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputIoTHub_updated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputIoTHubExists(resourceName),
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

func TestAccAzureRMStreamAnalyticsStreamInputIoTHub_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_stream_analytics_stream_input_iothub.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsStreamInputIoTHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsStreamInputIoTHub_json(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsStreamInputIoTHubExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMStreamAnalyticsStreamInputIoTHub_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_stream_analytics_stream_input_iothub"),
			},
		},
	})
}

func testCheckAzureRMStreamAnalyticsStreamInputIoTHubExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		jobName := rs.Primary.Attributes["stream_analytics_job_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := testAccProvider.Meta().(*ArmClient).streamanalytics.InputsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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
	conn := testAccProvider.Meta().(*ArmClient).streamanalytics.InputsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_stream_analytics_stream_input_iothub" {
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
			return fmt.Errorf("Stream Analytics Stream Input IoTHub still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMStreamAnalyticsStreamInputIoTHub_avro(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsStreamInputIoTHub_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_stream_input_iothub" "test" {
  name                         = "acctestinput-%d"
  stream_analytics_job_name    = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name          = "${azurerm_stream_analytics_job.test.resource_group_name}"
  endpoint                     = "messages/events"
  iothub_namespace             = "${azurerm_iothub.test.name}"
  eventhub_consumer_group_name = "$Default"
  shared_access_policy_key     = "${azurerm_iothub.test.shared_access_policy.0.primary_key}"
  shared_access_policy_name    = "iothubowner"

  serialization {
    type = "Avro"
  }
}
`, template, rInt)
}

func testAccAzureRMStreamAnalyticsStreamInputIoTHub_csv(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsStreamInputIoTHub_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_stream_input_iothub" "test" {
  name                         = "acctestinput-%d"
  stream_analytics_job_name    = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name          = "${azurerm_stream_analytics_job.test.resource_group_name}"
  endpoint                     = "messages/events"
  iothub_namespace             = "${azurerm_iothub.test.name}"
  eventhub_consumer_group_name = "$Default"
  shared_access_policy_key     = "${azurerm_iothub.test.shared_access_policy.0.primary_key}"
  shared_access_policy_name    = "iothubowner"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}
`, template, rInt)
}

func testAccAzureRMStreamAnalyticsStreamInputIoTHub_json(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsStreamInputIoTHub_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_stream_input_iothub" "test" {
  name                         = "acctestinput-%d"
  stream_analytics_job_name    = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name          = "${azurerm_stream_analytics_job.test.resource_group_name}"
  endpoint                     = "messages/events"
  iothub_namespace             = "${azurerm_iothub.test.name}"
  eventhub_consumer_group_name = "$Default"
  shared_access_policy_key     = "${azurerm_iothub.test.shared_access_policy.0.primary_key}"
  shared_access_policy_name    = "iothubowner"

  serialization {
    type     = "Json"
    encoding = "UTF8"
  }
}
`, template, rInt)
}

func testAccAzureRMStreamAnalyticsStreamInputIoTHub_updated(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsStreamInputIoTHub_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub" "updated" {
  name                = "acctestiot2-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }
}

resource "azurerm_stream_analytics_stream_input_iothub" "test" {
  name                         = "acctestinput-%d"
  stream_analytics_job_name    = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name          = "${azurerm_stream_analytics_job.test.resource_group_name}"
  endpoint                     = "messages/events"
  eventhub_consumer_group_name = "$Default"
  iothub_namespace             = "${azurerm_iothub.updated.name}"
  shared_access_policy_key     = "${azurerm_iothub.updated.shared_access_policy.0.primary_key}"
  shared_access_policy_name    = "iothubowner"

  serialization {
    type = "Avro"
  }
}
`, template, rInt, rInt)
}

func testAccAzureRMStreamAnalyticsStreamInputIoTHub_requiresImport(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsStreamInputIoTHub_json(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_stream_input_eventhub" "import" {
  name                         = "${azurerm_stream_analytics_stream_input_eventhub.test.name}"
  stream_analytics_job_name    = "${azurerm_stream_analytics_stream_input_eventhub.test.stream_analytics_job_name}"
  resource_group_name          = "${azurerm_stream_analytics_stream_input_eventhub.test.resource_group_name}"
  endpoint                     = "${azurerm_stream_analytics_stream_input_eventhub.test.endpoint}"
  eventhub_consumer_group_name = "${azurerm_stream_analytics_stream_input_eventhub.test.eventhub_consumer_group_name}"
  iothub_namespace             = "${azurerm_stream_analytics_stream_input_eventhub.test.iothub_namespace}"
  shared_access_policy_key     = "${azurerm_stream_analytics_stream_input_eventhub.test.shared_access_policy_key}"
  shared_access_policy_name    = "${azurerm_stream_analytics_stream_input_eventhub.test.resource_group_name}"
  serialization                = "${azurerm_stream_analytics_stream_input_eventhub.test.serialization}"
}
`, template)
}

func testAccAzureRMStreamAnalyticsStreamInputIoTHub_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestiothub-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }
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
`, rInt, location, rInt, rInt)
}
