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

func TestAccAzureRMStreamAnalyticsReferenceInputBlob_avro(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_reference_input_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsReferenceInputBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsReferenceInputBlob_avro(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsReferenceInputBlobExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_key"),
		},
	})
}

func TestAccAzureRMStreamAnalyticsReferenceInputBlob_csv(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_reference_input_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsReferenceInputBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsReferenceInputBlob_csv(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsReferenceInputBlobExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_key"),
		},
	})
}

func TestAccAzureRMStreamAnalyticsReferenceInputBlob_json(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_reference_input_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsReferenceInputBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsReferenceInputBlob_json(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsReferenceInputBlobExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_key"),
		},
	})
}

func TestAccAzureRMStreamAnalyticsReferenceInputBlob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_reference_input_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsReferenceInputBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsReferenceInputBlob_json(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsReferenceInputBlobExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMStreamAnalyticsReferenceInputBlob_updated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsReferenceInputBlobExists(data.ResourceName),
				),
			},
			data.ImportStep("storage_account_key"),
		},
	})
}

func TestAccAzureRMStreamAnalyticsReferenceInputBlob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_reference_input_blob", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsReferenceInputBlobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsReferenceInputBlob_json(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsReferenceInputBlobExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStreamAnalyticsReferenceInputBlob_requiresImport),
		},
	})
}

func testCheckAzureRMStreamAnalyticsReferenceInputBlobExists(resourceName string) resource.TestCheckFunc {
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

func testCheckAzureRMStreamAnalyticsReferenceInputBlobDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).StreamAnalytics.InputsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_stream_analytics_reference_input_blob" {
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
			return fmt.Errorf("Stream Analytics Stream Input Blob still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testAccAzureRMStreamAnalyticsReferenceInputBlob_avro(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsReferenceInputBlob_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_reference_input_blob" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  storage_account_name      = azurerm_storage_account.test.name
  storage_account_key       = azurerm_storage_account.test.primary_access_key
  storage_container_name    = azurerm_storage_container.test.name
  path_pattern              = "some-random-pattern"
  date_format               = "yyyy/MM/dd"
  time_format               = "HH"

  serialization {
    type = "Avro"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMStreamAnalyticsReferenceInputBlob_csv(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsReferenceInputBlob_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_reference_input_blob" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  storage_account_name      = azurerm_storage_account.test.name
  storage_account_key       = azurerm_storage_account.test.primary_access_key
  storage_container_name    = azurerm_storage_container.test.name
  path_pattern              = "some-random-pattern"
  date_format               = "yyyy/MM/dd"
  time_format               = "HH"

  serialization {
    type            = "Csv"
    encoding        = "UTF8"
    field_delimiter = ","
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMStreamAnalyticsReferenceInputBlob_json(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsReferenceInputBlob_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_reference_input_blob" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name
  storage_account_name      = azurerm_storage_account.test.name
  storage_account_key       = azurerm_storage_account.test.primary_access_key
  storage_container_name    = azurerm_storage_container.test.name
  path_pattern              = "some-random-pattern"
  date_format               = "yyyy/MM/dd"
  time_format               = "HH"

  serialization {
    type     = "Json"
    encoding = "UTF8"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMStreamAnalyticsReferenceInputBlob_updated(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsReferenceInputBlob_template(data)
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
  name                  = "example2"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_stream_analytics_reference_input_blob" "test" {
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

func testAccAzureRMStreamAnalyticsReferenceInputBlob_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsReferenceInputBlob_json(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_reference_input_blob" "import" {
  name                      = azurerm_stream_analytics_reference_input_blob.test.name
  stream_analytics_job_name = azurerm_stream_analytics_reference_input_blob.test.stream_analytics_job_name
  resource_group_name       = azurerm_stream_analytics_reference_input_blob.test.resource_group_name
  storage_account_name      = azurerm_stream_analytics_reference_input_blob.test.storage_account_name
  storage_account_key       = azurerm_stream_analytics_reference_input_blob.test.storage_account_key
  storage_container_name    = azurerm_stream_analytics_reference_input_blob.test.storage_container_name
  path_pattern              = azurerm_stream_analytics_reference_input_blob.test.path_pattern
  date_format               = azurerm_stream_analytics_reference_input_blob.test.date_format
  time_format               = azurerm_stream_analytics_reference_input_blob.test.time_format
  dynamic "serialization" {
    for_each = azurerm_stream_analytics_reference_input_blob.test.serialization
    content {
      encoding = lookup(serialization.value, "encoding", null)
      type     = serialization.value.type
    }
  }
}
`, template)
}

func testAccAzureRMStreamAnalyticsReferenceInputBlob_template(data acceptance.TestData) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
