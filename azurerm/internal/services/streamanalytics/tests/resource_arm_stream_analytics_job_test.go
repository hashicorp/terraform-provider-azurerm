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

func TestAccAzureRMStreamAnalyticsJob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsJob_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStreamAnalyticsJob_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsJob_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsJobExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStreamAnalyticsJob_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsJob_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsJobExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStreamAnalyticsJob_requiresImport),
		},
	})
}

func TestAccAzureRMStreamAnalyticsJob_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsJob_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsJobExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMStreamAnalyticsJob_updated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsJobExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMStreamAnalyticsJobExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).StreamAnalytics.JobsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on streamAnalyticsJobsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Stream Analytics Job %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMStreamAnalyticsJobDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).StreamAnalytics.JobsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_stream_analytics_job" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		resp, err := conn.Get(ctx, resourceGroup, name, "")
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Stream Analytics Job still exists:\n%#v", resp.StreamingJobProperties)
		}
	}

	return nil
}

func testAccAzureRMStreamAnalyticsJob_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_stream_analytics_job" "test" {
  name                = "acctestjob-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  streaming_units     = 3

  tags = {
    environment = "Test"
  }

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMStreamAnalyticsJob_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_stream_analytics_job" "test" {
  name                                     = "acctestjob-%d"
  resource_group_name                      = azurerm_resource_group.test.name
  location                                 = azurerm_resource_group.test.location
  data_locale                              = "en-GB"
  compatibility_level                      = "1.0"
  events_late_arrival_max_delay_in_seconds = 60
  events_out_of_order_max_delay_in_seconds = 50
  events_out_of_order_policy               = "Adjust"
  output_error_policy                      = "Drop"
  streaming_units                          = 3

  tags = {
    environment = "Test"
  }

  transformation_query = <<QUERY
    SELECT *
    INTO [YourOutputAlias]
    FROM [YourInputAlias]
QUERY

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMStreamAnalyticsJob_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsJob_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_job" "import" {
  name                                     = azurerm_stream_analytics_job.test.name
  resource_group_name                      = azurerm_stream_analytics_job.test.resource_group_name
  location                                 = azurerm_stream_analytics_job.test.location
  compatibility_level                      = azurerm_stream_analytics_job.test.compatibility_level
  data_locale                              = azurerm_stream_analytics_job.test.data_locale
  events_late_arrival_max_delay_in_seconds = azurerm_stream_analytics_job.test.events_late_arrival_max_delay_in_seconds
  events_out_of_order_max_delay_in_seconds = azurerm_stream_analytics_job.test.events_out_of_order_max_delay_in_seconds
  events_out_of_order_policy               = azurerm_stream_analytics_job.test.events_out_of_order_policy
  output_error_policy                      = azurerm_stream_analytics_job.test.output_error_policy
  streaming_units                          = azurerm_stream_analytics_job.test.streaming_units
  transformation_query                     = azurerm_stream_analytics_job.test.transformation_query
  tags                                     = azurerm_stream_analytics_job.test.tags
}
`, template)
}

func testAccAzureRMStreamAnalyticsJob_updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_stream_analytics_job" "test" {
  name                                     = "acctestjob-%d"
  resource_group_name                      = azurerm_resource_group.test.name
  location                                 = azurerm_resource_group.test.location
  data_locale                              = "en-GB"
  compatibility_level                      = "1.1"
  events_late_arrival_max_delay_in_seconds = 10
  events_out_of_order_max_delay_in_seconds = 20
  events_out_of_order_policy               = "Drop"
  output_error_policy                      = "Stop"
  streaming_units                          = 6

  transformation_query = <<QUERY
    SELECT *
    INTO [SomeOtherOutputAlias]
    FROM [SomeOtherInputAlias]
QUERY

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
