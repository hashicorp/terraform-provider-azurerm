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

func TestAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_function_javascript_udf", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_function_javascript_udf", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_requiresImport),
		},
	})
}

func TestAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_inputs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_function_javascript_udf", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_inputs(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).StreamAnalytics.FunctionsClient
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
			return fmt.Errorf("Bad: Get on streamAnalyticsFunctionsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Function JavaScript UDF %q (Stream Analytics Job %q / Resource Group %q) does not exist", name, jobName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).StreamAnalytics.OutputsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_stream_analytics_function_javascript_udf" {
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
			return fmt.Errorf("Stream Analytics Function JavaScript UDF still exists:\n%#v", resp.OutputProperties)
		}
	}

	return nil
}

func testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_basic(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_function_javascript_udf" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name

  script = <<SCRIPT
function getRandomNumber(in) {
  return in;
}
SCRIPT


  input {
    type = "bigint"
  }

  output {
    type = "bigint"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_function_javascript_udf" "import" {
  name                      = azurerm_stream_analytics_function_javascript_udf.test.name
  stream_analytics_job_name = azurerm_stream_analytics_function_javascript_udf.test.stream_analytics_job_name
  resource_group_name       = azurerm_stream_analytics_function_javascript_udf.test.resource_group_name
  script                    = azurerm_stream_analytics_function_javascript_udf.test.script

  input {
    type = azurerm_stream_analytics_function_javascript_udf.test.input.0.type
  }

  output {
    type = azurerm_stream_analytics_function_javascript_udf.test.output.0.type
  }
}
`, template)
}

func testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_inputs(data acceptance.TestData) string {
	template := testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_function_javascript_udf" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = azurerm_stream_analytics_job.test.name
  resource_group_name       = azurerm_stream_analytics_job.test.resource_group_name

  script = <<SCRIPT
function getRandomNumber(first, second) {
  return first * second;
}
SCRIPT


  input {
    type = "bigint"
  }

  input {
    type = "bigint"
  }

  output {
    type = "bigint"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
