package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMStreamAnalyticsOutputEventHub_avro(t *testing.T) {
	resourceName := "azurerm_stream_analytics_output_eventhub.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputEventHub_avro(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputEventHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "serialization.0.type", "Avro"),
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

func TestAccAzureRMStreamAnalyticsOutputEventHub_csv(t *testing.T) {
	resourceName := "azurerm_stream_analytics_output_eventhub.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputEventHub_csv(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputEventHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "serialization.0.type", "Csv"),
					resource.TestCheckResourceAttr(resourceName, "serialization.0.field_delimiter", ","),
					resource.TestCheckResourceAttr(resourceName, "serialization.0.encoding", "UTF8"),
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

func TestAccAzureRMStreamAnalyticsOutputEventHub_json(t *testing.T) {
	resourceName := "azurerm_stream_analytics_output_eventhub.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputEventHub_json(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputEventHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "serialization.0.format", "LineSeparated"),
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

func TestAccAzureRMStreamAnalyticsOutputEventHub_jsonArrayFormat(t *testing.T) {
	resourceName := "azurerm_stream_analytics_output_eventhub.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputEventHub_jsonArrayFormat(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputEventHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "serialization.0.format", "Array"),
					resource.TestCheckResourceAttr(resourceName, "serialization.0.type", "Json"),
					resource.TestCheckResourceAttr(resourceName, "serialization.0.encoding", "UTF8"),
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

func TestAccAzureRMStreamAnalyticsOutputEventHub_update(t *testing.T) {
	resourceName := "azurerm_stream_analytics_output_eventhub.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputEventHub_json(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputEventHubExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMStreamAnalyticsOutputEventHub_updated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputEventHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "serialization.0.type", "Avro"),
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

func TestAccAzureRMStreamAnalyticsOutputEventHub_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_stream_analytics_output_eventhub.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputEventHub_json(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputEventHubExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMStreamAnalyticsOutputEventHub_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_stream_analytics_output_eventhub"),
			},
		},
	})
}

func testCheckAzureRMStreamAnalyticsOutputEventHubExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		jobName := rs.Primary.Attributes["stream_analytics_job_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := acceptance.AzureProvider.Meta().(*clients.Client).StreamAnalytics.OutputsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
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

func testCheckAzureRMStreamAnalyticsOutputEventHubDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).StreamAnalytics.OutputsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_stream_analytics_output_eventhub" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		jobName := rs.Primary.Attributes["stream_analytics_job_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
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

func testAccAzureRMStreamAnalyticsOutputEventHub_avro(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputEventHub_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_eventhub" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name       = "${azurerm_stream_analytics_job.test.resource_group_name}"
  eventhub_name             = "${azurerm_eventhub.test.name}"
  servicebus_namespace      = "${azurerm_eventhub_namespace.test.name}"
  shared_access_policy_key  = "${azurerm_eventhub_namespace.test.default_primary_key}"
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type = "Avro"
  }
}
`, template, rInt)
}

func testAccAzureRMStreamAnalyticsOutputEventHub_csv(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputEventHub_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_eventhub" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name       = "${azurerm_stream_analytics_job.test.resource_group_name}"
  eventhub_name             = "${azurerm_eventhub.test.name}"
  servicebus_namespace      = "${azurerm_eventhub_namespace.test.name}"
  shared_access_policy_key  = "${azurerm_eventhub_namespace.test.default_primary_key}"
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}
`, template, rInt)
}

func testAccAzureRMStreamAnalyticsOutputEventHub_json(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputEventHub_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_eventhub" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name       = "${azurerm_stream_analytics_job.test.resource_group_name}"
  eventhub_name             = "${azurerm_eventhub.test.name}"
  servicebus_namespace      = "${azurerm_eventhub_namespace.test.name}"
  shared_access_policy_key  = "${azurerm_eventhub_namespace.test.default_primary_key}"
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type     = "Json"
    encoding = "UTF8"
    format   = "LineSeparated"
  }
}
`, template, rInt)
}

func testAccAzureRMStreamAnalyticsOutputEventHub_jsonArrayFormat(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputEventHub_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_eventhub" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name       = "${azurerm_stream_analytics_job.test.resource_group_name}"
  eventhub_name             = "${azurerm_eventhub.test.name}"
  servicebus_namespace      = "${azurerm_eventhub_namespace.test.name}"
  shared_access_policy_key  = "${azurerm_eventhub_namespace.test.default_primary_key}"
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type     = "Json"
    encoding = "UTF8"
    format   = "Array"
  }
}
`, template, rInt)
}

func testAccAzureRMStreamAnalyticsOutputEventHub_updated(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputEventHub_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_namespace" "updated" {
  name                = "acctestehn2-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  capacity            = 1
}

resource "azurerm_eventhub" "updated" {
  name                = "acctesteh2-%d"
  namespace_name      = "${azurerm_eventhub_namespace.updated.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_stream_analytics_output_eventhub" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name       = "${azurerm_stream_analytics_job.test.resource_group_name}"
  eventhub_name             = "${azurerm_eventhub.updated.name}"
  servicebus_namespace      = "${azurerm_eventhub_namespace.updated.name}"
  shared_access_policy_key  = "${azurerm_eventhub_namespace.updated.default_primary_key}"
  shared_access_policy_name = "RootManageSharedAccessKey"

  serialization {
    type = "Avro"
  }
}
`, template, rInt, rInt, rInt)
}

func testAccAzureRMStreamAnalyticsOutputEventHub_requiresImport(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputEventHub_json(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_eventhub" "import" {
  name                      = "${azurerm_stream_analytics_output_eventhub.test.name}"
  stream_analytics_job_name = "${azurerm_stream_analytics_output_eventhub.test.stream_analytics_job_name}"
  resource_group_name       = "${azurerm_stream_analytics_output_eventhub.test.resource_group_name}"
  eventhub_name             = "${azurerm_stream_analytics_output_eventhub.test.eventhub_name}"
  servicebus_namespace      = "${azurerm_stream_analytics_output_eventhub.test.servicebus_namespace}"
  shared_access_policy_key  = "${azurerm_stream_analytics_output_eventhub.test.shared_access_policy_key}"
  shared_access_policy_name = "${azurerm_stream_analytics_output_eventhub.test.shared_access_policy_name}"
}
`, template)
}

func testAccAzureRMStreamAnalyticsOutputEventHub_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctestehn-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
  capacity            = 1
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteh-%d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  partition_count     = 2
  message_retention   = 1
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
