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

func TestAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_basic(t *testing.T) {
	resourceName := "azurerm_stream_analytics_function_javascript_udf.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_stream_analytics_function_javascript_udf.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_stream_analytics_function_javascript_udf"),
			},
		},
	})
}

func TestAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_inputs(t *testing.T) {
	resourceName := "azurerm_stream_analytics_function_javascript_udf.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_inputs(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMStreamAnalyticsFunctionJavaScriptUDFExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		jobName := rs.Primary.Attributes["stream_analytics_job_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := testAccProvider.Meta().(*ArmClient).streamanalytics.FunctionsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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
	conn := testAccProvider.Meta().(*ArmClient).streamanalytics.OutputsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_stream_analytics_function_javascript_udf" {
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
			return fmt.Errorf("Stream Analytics Function JavaScript UDF still exists:\n%#v", resp.OutputProperties)
		}
	}

	return nil
}

func testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_basic(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_function_javascript_udf" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name       = "${azurerm_stream_analytics_job.test.resource_group_name}"

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
`, template, rInt)
}

func testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_requiresImport(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_function_javascript_udf" "import" {
  name                      = "${azurerm_stream_analytics_function_javascript_udf.test.name}"
  stream_analytics_job_name = "${azurerm_stream_analytics_function_javascript_udf.test.stream_analytics_job_name}"
  resource_group_name       = "${azurerm_stream_analytics_function_javascript_udf.test.resource_group_name}"
  script                    = "${azurerm_stream_analytics_function_javascript_udf.test.script}"
  inputs                    = "${azurerm_stream_analytics_function_javascript_udf.test.inputs}"
  output                    = "${azurerm_stream_analytics_function_javascript_udf.test.output}"
}
`, template)
}

func testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_inputs(rInt int, location string) string {
	template := testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_function_javascript_udf" "test" {
  name                      = "acctestinput-%d"
  stream_analytics_job_name = "${azurerm_stream_analytics_job.test.name}"
  resource_group_name       = "${azurerm_stream_analytics_job.test.resource_group_name}"

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
`, template, rInt)
}

func testAccAzureRMStreamAnalyticsFunctionJavaScriptUDF_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
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
`, rInt, location, rInt)
}
