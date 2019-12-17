package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMStreamAnalyticsOutputBlob_avro(t *testing.T) {
	resourceName := "azurerm_stream_analytics_output_blob.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(5)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputBlob_avro(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputBlobExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"storage_account_key",
				},
			},
		},
	})
}

func TestAccAzureRMStreamAnalyticsOutputBlob_csv(t *testing.T) {
	resourceName := "azurerm_stream_analytics_output_blob.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(5)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputBlob_csv(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputBlobExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"storage_account_key",
				},
			},
		},
	})
}

func TestAccAzureRMStreamAnalyticsOutputBlob_json(t *testing.T) {
	resourceName := "azurerm_stream_analytics_output_blob.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(5)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputBlob_json(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputBlobExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"storage_account_key",
				},
			},
		},
	})
}

func TestAccAzureRMStreamAnalyticsOutputBlob_update(t *testing.T) {
	resourceName := "azurerm_stream_analytics_output_blob.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(5)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputBlob_json(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputBlobExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMStreamAnalyticsOutputBlob_updated(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputBlobExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"storage_account_key",
				},
			},
		},
	})
}

func TestAccAzureRMStreamAnalyticsOutputBlob_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_stream_analytics_output_blob.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(5)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsOutputBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsOutputBlob_json(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsOutputBlobExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMStreamAnalyticsOutputBlob_requiresImport(ri, rs, location),
				ExpectError: acceptance.RequiresImportError("azurerm_stream_analytics_output_blob"),
			},
		},
	})
}

func testCheckAzureRMStreamAnalyticsOutputBlobExists(resourceName string) resource.TestCheckFunc {
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

func testCheckAzureRMStreamAnalyticsOutputBlobDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).StreamAnalytics.OutputsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_stream_analytics_output_blob" {
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

func testAccAzureRMStreamAnalyticsOutputBlob_avro(rInt int, rString string, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputBlob_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_blob" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name       = "${azurerm_stream_analytics_job.test.resource_group_name}"
  storage_account_name      = "${azurerm_storage_account.test.name}"
  storage_account_key       = "${azurerm_storage_account.test.primary_access_key}"
  storage_container_name    = "${azurerm_storage_container.test.name}"
  path_pattern              = "some-other-pattern"
  date_format               = "yyyy-MM-dd"
  time_format               = "HH"

  serialization {
    type = "Avro"
  }
}
`, template, rInt)
}

func testAccAzureRMStreamAnalyticsOutputBlob_csv(rInt int, rString string, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputBlob_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_blob" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name       = "${azurerm_stream_analytics_job.test.resource_group_name}"
  storage_account_name      = "${azurerm_storage_account.test.name}"
  storage_account_key       = "${azurerm_storage_account.test.primary_access_key}"
  storage_container_name    = "${azurerm_storage_container.test.name}"
  path_pattern              = "some-pattern"
  date_format               = "yyyy-MM-dd"
  time_format               = "HH"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}
`, template, rInt)
}

func testAccAzureRMStreamAnalyticsOutputBlob_json(rInt int, rString string, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputBlob_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_blob" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name       = "${azurerm_stream_analytics_job.test.resource_group_name}"
  storage_account_name      = "${azurerm_storage_account.test.name}"
  storage_account_key       = "${azurerm_storage_account.test.primary_access_key}"
  storage_container_name    = "${azurerm_storage_container.test.name}"
  path_pattern              = "some-pattern"
  date_format               = "yyyy-MM-dd"
  time_format               = "HH"

  serialization {
    type     = "Json"
    encoding = "UTF8"
    format   = "LineSeparated"
  }
}
`, template, rInt)
}

func testAccAzureRMStreamAnalyticsOutputBlob_updated(rInt int, rString string, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputBlob_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "updated" {
  name                     = "acctestsa2%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "updated" {
  name                  = "example"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.updated.name}"
  container_access_type = "private"
}

resource "azurerm_stream_analytics_output_blob" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name       = "${azurerm_stream_analytics_job.test.resource_group_name}"
  storage_account_name      = "${azurerm_storage_account.updated.name}"
  storage_account_key       = "${azurerm_storage_account.updated.primary_access_key}"
  storage_container_name    = "${azurerm_storage_container.updated.name}"
  path_pattern              = "some-other-pattern"
  date_format               = "yyyy-MM-dd"
  time_format               = "HH"

  serialization {
    type = "Avro"
  }
}
`, template, rString, rInt)
}

func testAccAzureRMStreamAnalyticsOutputBlob_requiresImport(rInt int, rString string, location string) string {
	template := testAccAzureRMStreamAnalyticsOutputBlob_json(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_output_blob" "import" {
  name                      = "${azurerm_stream_analytics_output_blob.test.name}"
  stream_analytics_job_name = "${azurerm_stream_analytics_output_blob.test.stream_analytics_job_name}"
  resource_group_name       = "${azurerm_stream_analytics_output_blob.test.resource_group_name}"
  storage_account_name      = "${azurerm_stream_analytics_output_blob.test.storage_account_name}"
  storage_account_key       = "${azurerm_stream_analytics_output_blob.test.storage_account_key}"
  storage_container_name    = "${azurerm_stream_analytics_output_blob.test.storage_container_name}"
  path_pattern              = "${azurerm_stream_analytics_output_blob.test.path_pattern}"
  date_format               = "${azurerm_stream_analytics_output_blob.test.date_format}"
  time_format               = "${azurerm_stream_analytics_output_blob.test.time_format}"
  serialization             = "${azurerm_stream_analytics_output_blob.test.serialization}"
}
`, template)
}

func testAccAzureRMStreamAnalyticsOutputBlob_template(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "example"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
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
`, rInt, location, rString, rInt)
}
